package ngap

import (
	"github.com/free5gc/ngap/ngapType"
)

func getCauseRadioNetworkErrorStr(cause *ngapType.CauseRadioNetwork) string {
	switch cause.Value {
	case ngapType.CauseRadioNetworkPresentUnspecified:
		return "RadioNetwork : Unspecified"
	case ngapType.CauseRadioNetworkPresentTxnrelocoverallExpiry:
		return "RadioNetwork : TxnrelocoverallExpiry"
	case ngapType.CauseRadioNetworkPresentSuccessfulHandover:
		return "RadioNetwork : SuccessfulHandover"
	case ngapType.CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason:
		return "RadioNetwork : ReleaseDueToNgranGeneratedReason"
	case ngapType.CauseRadioNetworkPresentReleaseDueTo5gcGeneratedReason:
		return "RadioNetwork : ReleaseDueTo5gcGeneratedReason"
	case ngapType.CauseRadioNetworkPresentHandoverCancelled:
		return "RadioNetwork : HandoverCancelled"
	case ngapType.CauseRadioNetworkPresentPartialHandover:
		return "RadioNetwork : PartialHandover"
	case ngapType.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem:
		return "RadioNetwork : HoFailureInTarget5GCNgranNodeOrTargetSystem"
	case ngapType.CauseRadioNetworkPresentHoTargetNotAllowed:
		return "RadioNetwork : HoTargetNotAllowed"
	case ngapType.CauseRadioNetworkPresentTngrelocoverallExpiry:
		return "RadioNetwork : TngrelocoverallExpiry"
	case ngapType.CauseRadioNetworkPresentTngrelocprepExpiry:
		return "RadioNetwork : TngrelocprepExpiry"
	case ngapType.CauseRadioNetworkPresentCellNotAvailable:
		return "RadioNetwork : CellNotAvailable"
	case ngapType.CauseRadioNetworkPresentUnknownTargetID:
		return "RadioNetwork : UnknownTargetID"
	case ngapType.CauseRadioNetworkPresentNoRadioResourcesAvailableInTargetCell:
		return "RadioNetwork : NoRadioResourcesAvailableInTargetCell"
	case ngapType.CauseRadioNetworkPresentUnknownLocalUENGAPID:
		return "RadioNetwork : UnknownLocalUENGAPID"
	case ngapType.CauseRadioNetworkPresentInconsistentRemoteUENGAPID:
		return "RadioNetwork : InconsistentRemoteUENGAPID"
	case ngapType.CauseRadioNetworkPresentHandoverDesirableForRadioReason:
		return "RadioNetwork : HandoverDesirableForRadioReason"
	case ngapType.CauseRadioNetworkPresentTimeCriticalHandover:
		return "RadioNetwork : TimeCriticalHandover"
	case ngapType.CauseRadioNetworkPresentResourceOptimisationHandover:
		return "RadioNetwork : ResourceOptimisationHandover"
	case ngapType.CauseRadioNetworkPresentReduceLoadInServingCell:
		return "RadioNetwork : ReduceLoadInServingCell"
	case ngapType.CauseRadioNetworkPresentUserInactivity:
		return "RadioNetwork : UserInactivity"
	case ngapType.CauseRadioNetworkPresentRadioConnectionWithUeLost:
		return "RadioNetwork : RadioConnectionWithUeLost"
	case ngapType.CauseRadioNetworkPresentRadioResourcesNotAvailable:
		return "RadioNetwork : RadioResourcesNotAvailable"
	case ngapType.CauseRadioNetworkPresentInvalidQosCombination:
		return "RadioNetwork : InvalidQosCombination"
	case ngapType.CauseRadioNetworkPresentFailureInRadioInterfaceProcedure:
		return "RadioNetwork : FailureInRadioInterfaceProcedure"
	case ngapType.CauseRadioNetworkPresentInteractionWithOtherProcedure:
		return "RadioNetwork : InteractionWithOtherProcedure"
	case ngapType.CauseRadioNetworkPresentUnknownPDUSessionID:
		return "RadioNetwork : UnknownPDUSessionID"
	case ngapType.CauseRadioNetworkPresentUnkownQosFlowID:
		return "RadioNetwork : UnkownQosFlowID"
	case ngapType.CauseRadioNetworkPresentMultiplePDUSessionIDInstances:
		return "RadioNetwork : MultiplePDUSessionIDInstances"
	case ngapType.CauseRadioNetworkPresentMultipleQosFlowIDInstances:
		return "RadioNetwork : MultipleQosFlowIDInstances"
	case ngapType.CauseRadioNetworkPresentEncryptionAndOrIntegrityProtectionAlgorithmsNotSupported:
		return "RadioNetwork : EncryptionAndOrIntegrityProtectionAlgorithmsNotSupported"
	case ngapType.CauseRadioNetworkPresentNgIntraSystemHandoverTriggered:
		return "RadioNetwork : NgIntraSystemHandoverTriggered"
	case ngapType.CauseRadioNetworkPresentNgInterSystemHandoverTriggered:
		return "RadioNetwork : NgInterSystemHandoverTriggered"
	case ngapType.CauseRadioNetworkPresentXnHandoverTriggered:
		return "RadioNetwork : XnHandoverTriggered"
	case ngapType.CauseRadioNetworkPresentNotSupported5QIValue:
		return "RadioNetwork : NotSupported5QIValue"
	case ngapType.CauseRadioNetworkPresentUeContextTransfer:
		return "RadioNetwork : UeContextTransfer"
	case ngapType.CauseRadioNetworkPresentImsVoiceEpsFallbackOrRatFallbackTriggered:
		return "RadioNetwork : ImsVoiceEpsFallbackOrRatFallbackTriggered"
	case ngapType.CauseRadioNetworkPresentUpIntegrityProtectionNotPossible:
		return "RadioNetwork : UpIntegrityProtectionNotPossible"
	case ngapType.CauseRadioNetworkPresentUpConfidentialityProtectionNotPossible:
		return "RadioNetwork : UpConfidentialityProtectionNotPossible"
	case ngapType.CauseRadioNetworkPresentSliceNotSupported:
		return "RadioNetwork : SliceNotSupported"
	case ngapType.CauseRadioNetworkPresentUeInRrcInactiveStateNotReachable:
		return "RadioNetwork : UeInRrcInactiveStateNotReachable"
	case ngapType.CauseRadioNetworkPresentRedirection:
		return "RadioNetwork : Redirection"
	case ngapType.CauseRadioNetworkPresentResourcesNotAvailableForTheSlice:
		return "RadioNetwork : ResourcesNotAvailableForTheSlice"
	case ngapType.CauseRadioNetworkPresentUeMaxIntegrityProtectedDataRateReason:
		return "RadioNetwork : UeMaxIntegrityProtectedDataRateReason"
	case ngapType.CauseRadioNetworkPresentReleaseDueToCnDetectedMobility:
		return "RadioNetwork : ReleaseDueToCnDetectedMobility"
	case ngapType.CauseRadioNetworkPresentN26InterfaceNotAvailable:
		return "RadioNetwork : N26InterfaceNotAvailable"
	case ngapType.CauseRadioNetworkPresentReleaseDueToPreEmption:
		return "RadioNetwork : ReleaseDueToPreEmption"
	default:
		return "unknown cause"
	}
}

func getCauseTransportErrorStr(cause *ngapType.CauseTransport) string {
	switch cause.Value {
	case ngapType.CauseTransportPresentTransportResourceUnavailable:
		return "Transport : TransportResourceUnavailable"
	case ngapType.CauseTransportPresentUnspecified:
		return "Transport : Unspecified"
	default:
		return "unknown cause"
	}
}

func getCauseNasErrorStr(cause *ngapType.CauseNas) string {
	switch cause.Value {
	case ngapType.CauseNasPresentNormalRelease:
		return "Nas : NormalRelease"
	case ngapType.CauseNasPresentAuthenticationFailure:
		return "Nas : AuthenticationFailure"
	case ngapType.CauseNasPresentDeregister:
		return "Nas : Deregister"
	case ngapType.CauseNasPresentUnspecified:
		return "Nas : Unspecified"
	default:
		return "unknown cause"
	}
}

func getCauseProtocolErrorStr(cause *ngapType.CauseProtocol) string {
	switch cause.Value {
	case ngapType.CauseProtocolPresentTransferSyntaxError:
		return "Protocol : TransferSyntaxError"
	case ngapType.CauseProtocolPresentAbstractSyntaxErrorReject:
		return "Protocol : AbstractSyntaxErrorReject"
	case ngapType.CauseProtocolPresentAbstractSyntaxErrorIgnoreAndNotify:
		return "Protocol : AbstractSyntaxErrorIgnoreAndNotify"
	case ngapType.CauseProtocolPresentMessageNotCompatibleWithReceiverState:
		return "Protocol : MessageNotCompatibleWithReceiverState"
	case ngapType.CauseProtocolPresentSemanticError:
		return "Protocol : SemanticError"
	case ngapType.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage:
		return "Protocol : AbstractSyntaxErrorFalselyConstructedMessage"
	case ngapType.CauseProtocolPresentUnspecified:
		return "Protocol : Unspecified"
	default:
		return "unknown cause"
	}
}

func getCauseMiscErrorStr(cause *ngapType.CauseMisc) string {
	switch cause.Value {
	case ngapType.CauseMiscPresentControlProcessingOverload:
		return "Misc : ControlProcessingOverload"
	case ngapType.CauseMiscPresentNotEnoughUserPlaneProcessingResources:
		return "Misc : NotEnoughUserPlaneProcessingResources"
	case ngapType.CauseMiscPresentHardwareFailure:
		return "Misc : HardwareFailure"
	case ngapType.CauseMiscPresentOmIntervention:
		return "Misc : OmIntervention"
	case ngapType.CauseMiscPresentUnknownPLMN:
		return "Misc : UnknownPLMN"
	case ngapType.CauseMiscPresentUnspecified:
		return "Misc : Unspecified"
	default:
		return "unknown cause"
	}
}

func getCauseChoiceExtensionsErrorStr(cause *ngapType.ProtocolIESingleContainerCauseExtIEs) string {
	return "ChoiceExtensions : Unknown error"
}

func GetCauseErrorStr(cause *ngapType.Cause) string {
	if cause != nil {
		switch cause.Present {
		case ngapType.CausePresentRadioNetwork:
			return getCauseRadioNetworkErrorStr(cause.RadioNetwork)
		case ngapType.CausePresentTransport:
			return getCauseTransportErrorStr(cause.Transport)
		case ngapType.CausePresentNas:
			return getCauseNasErrorStr(cause.Nas)
		case ngapType.CausePresentProtocol:
			return getCauseProtocolErrorStr(cause.Protocol)
		case ngapType.CausePresentMisc:
			return getCauseMiscErrorStr(cause.Misc)
		case ngapType.CausePresentChoiceExtensions:
			return getCauseChoiceExtensionsErrorStr(cause.ChoiceExtensions)
		default:
			return UNKNOWN_NGAP_TYPE_CAUSE_ERR
		}
	}

	return UNKNOWN_NGAP_TYPE_CAUSE_ERR
}
