package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	clienttypes "github.com/cosmos/ibc-go/v5/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v5/modules/core/04-channel/types"
	host "github.com/cosmos/ibc-go/v5/modules/core/24-host"

	spntypes "github.com/tendermint/spn/pkg/types"
	"github.com/tendermint/spn/x/monitoringp/types"
)

// TransmitMonitoringPacket transmits the packet over IBC with the specified source port and source channel
func (k Keeper) TransmitMonitoringPacket(
	ctx sdk.Context,
	packetData spntypes.MonitoringPacket,
	sourcePort,
	sourceChannel string,
	timeoutHeight clienttypes.Height,
	timeoutTimestamp uint64,
) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(channeltypes.ErrChannelNotFound, "port ID (%s) channel ID (%s)", sourcePort, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrapf(
			channeltypes.ErrSequenceSendNotFound,
			"source port: %s, source channel: %s", sourcePort, sourceChannel,
		)
	}

	channelCap, ok := k.scopedKeeper.GetCapability(ctx, host.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}

	// encode the packet
	var modulePacket spntypes.MonitoringPacketData
	modulePacket.Packet = &spntypes.MonitoringPacketData_MonitoringPacket{
		MonitoringPacket: &packetData,
	}

	packetBytes, err := types.ModuleCdc.MarshalJSON(&modulePacket)
	if err != nil {
		return sdkerrors.Wrap(types.ErrJSONMarshal, err.Error())
	}

	packet := channeltypes.NewPacket(
		packetBytes,
		sequence,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		timeoutHeight,
		timeoutTimestamp,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

// OnRecvMonitoringPacket processes packet reception
func (k Keeper) OnRecvMonitoringPacket(
	_ sdk.Context,
	_ channeltypes.Packet,
	_ spntypes.MonitoringPacket,
) (packetAck spntypes.MonitoringPacketAck, err error) {
	return packetAck, types.ErrNotImplemented
}

// OnAcknowledgementMonitoringPacket responds to the the success or failure of a packet
// acknowledgement written on the receiving chain.
func (k Keeper) OnAcknowledgementMonitoringPacket(
	ctx sdk.Context,
	packet channeltypes.Packet,
	data spntypes.MonitoringPacket,
	ack channeltypes.Acknowledgement,
) error {
	switch dispatchedAck := ack.Response.(type) {
	case *channeltypes.Acknowledgement_Error:
		return nil
	case *channeltypes.Acknowledgement_Result:
		// Decode the packet acknowledgment
		var packetAck spntypes.MonitoringPacketAck

		if err := types.ModuleCdc.UnmarshalJSON(dispatchedAck.Result, &packetAck); err != nil {
			// The counter-party module doesn't implement the correct acknowledgment format
			return sdkerrors.Wrap(types.ErrJSONUnmarshal, err.Error())
		}

		return nil
	default:
		// The counter-party module doesn't implement the correct acknowledgment format
		return sdkerrors.Wrapf(types.ErrUnrecognizedAckType, "ack type: %T", ack)

	}
}

// OnTimeoutMonitoringPacket responds to the case where a packet has not been transmitted because of a timeout
func (k Keeper) OnTimeoutMonitoringPacket(
	_ sdk.Context,
	_ channeltypes.Packet,
	_ spntypes.MonitoringPacket,
) error {
	return types.ErrNotImplemented
}
