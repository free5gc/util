// Package ngap: NGAP message metrics, built against the new
// github.com/free5gc/ngap API (ie package).
package ngap

import (
	"github.com/free5gc/ngap/aper"
	ngapie "github.com/free5gc/ngap/ie"
)

// The cause maps below pin the "cause" label of the NGAP metrics to this
// package instead of delegating to ngapie.Cause.String(). Two reasons:
//
//   - The label values are consumed by dashboards and alerting rules, while the
//     strings of github.com/free5gc/ngap belong to that module and may change
//     with no compilation error here.
//   - ngapie.Cause.String() renders values it does not know as
//     "Unregconized <group> cause value(<n>)", which turns a malformed peer into
//     unbounded label cardinality. An unknown value maps to the fixed
//     UNKNOWN_CAUSE_ERR here.
//
// Every entry is "<group> : <constant name without the CauseXxxPresent prefix>",
// matching the values emitted before the migration to the new ngap API.
var causeRadioNetworkStr = map[aper.Enumerated]string{
	ngapie.CauseRadioNetworkPresentUnspecified:                      "RadioNetwork : Unspecified",
	ngapie.CauseRadioNetworkPresentTxnrelocoverallExpiry:            "RadioNetwork : TxnrelocoverallExpiry",
	ngapie.CauseRadioNetworkPresentSuccessfulHandover:               "RadioNetwork : SuccessfulHandover",
	ngapie.CauseRadioNetworkPresentReleaseDueToNgranGeneratedReason: "RadioNetwork : ReleaseDueToNgranGeneratedReason",
	ngapie.CauseRadioNetworkPresentReleaseDueTo5gcGeneratedReason:   "RadioNetwork : ReleaseDueTo5gcGeneratedReason",
	ngapie.CauseRadioNetworkPresentHandoverCancelled:                "RadioNetwork : HandoverCancelled",
	ngapie.CauseRadioNetworkPresentPartialHandover:                  "RadioNetwork : PartialHandover",
	ngapie.CauseRadioNetworkPresentHoFailureInTarget5GCNgranNodeOrTargetSystem: "RadioNetwork : " +
		"HoFailureInTarget5GCNgranNodeOrTargetSystem",
	ngapie.CauseRadioNetworkPresentHoTargetNotAllowed:    "RadioNetwork : HoTargetNotAllowed",
	ngapie.CauseRadioNetworkPresentTngrelocoverallExpiry: "RadioNetwork : TngrelocoverallExpiry",
	ngapie.CauseRadioNetworkPresentTngrelocprepExpiry:    "RadioNetwork : TngrelocprepExpiry",
	ngapie.CauseRadioNetworkPresentCellNotAvailable:      "RadioNetwork : CellNotAvailable",
	ngapie.CauseRadioNetworkPresentUnknownTargetID:       "RadioNetwork : UnknownTargetID",
	ngapie.CauseRadioNetworkPresentNoRadioResourcesAvailableInTargetCell: "RadioNetwork : " +
		"NoRadioResourcesAvailableInTargetCell",
	ngapie.CauseRadioNetworkPresentUnknownLocalUENGAPID:             "RadioNetwork : UnknownLocalUENGAPID",
	ngapie.CauseRadioNetworkPresentInconsistentRemoteUENGAPID:       "RadioNetwork : InconsistentRemoteUENGAPID",
	ngapie.CauseRadioNetworkPresentHandoverDesirableForRadioReason:  "RadioNetwork : HandoverDesirableForRadioReason",
	ngapie.CauseRadioNetworkPresentTimeCriticalHandover:             "RadioNetwork : TimeCriticalHandover",
	ngapie.CauseRadioNetworkPresentResourceOptimisationHandover:     "RadioNetwork : ResourceOptimisationHandover",
	ngapie.CauseRadioNetworkPresentReduceLoadInServingCell:          "RadioNetwork : ReduceLoadInServingCell",
	ngapie.CauseRadioNetworkPresentUserInactivity:                   "RadioNetwork : UserInactivity",
	ngapie.CauseRadioNetworkPresentRadioConnectionWithUeLost:        "RadioNetwork : RadioConnectionWithUeLost",
	ngapie.CauseRadioNetworkPresentRadioResourcesNotAvailable:       "RadioNetwork : RadioResourcesNotAvailable",
	ngapie.CauseRadioNetworkPresentInvalidQosCombination:            "RadioNetwork : InvalidQosCombination",
	ngapie.CauseRadioNetworkPresentFailureInRadioInterfaceProcedure: "RadioNetwork : FailureInRadioInterfaceProcedure",
	ngapie.CauseRadioNetworkPresentInteractionWithOtherProcedure:    "RadioNetwork : InteractionWithOtherProcedure",
	ngapie.CauseRadioNetworkPresentUnknownPDUSessionID:              "RadioNetwork : UnknownPDUSessionID",
	ngapie.CauseRadioNetworkPresentUnkownQosFlowID:                  "RadioNetwork : UnkownQosFlowID",
	ngapie.CauseRadioNetworkPresentMultiplePDUSessionIDInstances:    "RadioNetwork : MultiplePDUSessionIDInstances",
	ngapie.CauseRadioNetworkPresentMultipleQosFlowIDInstances:       "RadioNetwork : MultipleQosFlowIDInstances",
	ngapie.CauseRadioNetworkPresentEncryptionAndOrIntegrityProtectionAlgorithmsNotSupported: "RadioNetwork : " +
		"EncryptionAndOrIntegrityProtectionAlgorithmsNotSupported",
	ngapie.CauseRadioNetworkPresentNgIntraSystemHandoverTriggered: "RadioNetwork : NgIntraSystemHandoverTriggered",
	ngapie.CauseRadioNetworkPresentNgInterSystemHandoverTriggered: "RadioNetwork : NgInterSystemHandoverTriggered",
	ngapie.CauseRadioNetworkPresentXnHandoverTriggered:            "RadioNetwork : XnHandoverTriggered",
	ngapie.CauseRadioNetworkPresentNotSupported5QIValue:           "RadioNetwork : NotSupported5QIValue",
	ngapie.CauseRadioNetworkPresentUeContextTransfer:              "RadioNetwork : UeContextTransfer",
	ngapie.CauseRadioNetworkPresentImsVoiceEpsFallbackOrRatFallbackTriggered: "RadioNetwork : " +
		"ImsVoiceEpsFallbackOrRatFallbackTriggered",
	ngapie.CauseRadioNetworkPresentUpIntegrityProtectionNotPossible: "RadioNetwork : UpIntegrityProtectionNotPossible",
	ngapie.CauseRadioNetworkPresentUpConfidentialityProtectionNotPossible: "RadioNetwork : " +
		"UpConfidentialityProtectionNotPossible",
	ngapie.CauseRadioNetworkPresentSliceNotSupported: "RadioNetwork : SliceNotSupported",
	ngapie.CauseRadioNetworkPresentUeInRrcInactiveStateNotReachable: "RadioNetwork : " +
		"UeInRrcInactiveStateNotReachable",
	ngapie.CauseRadioNetworkPresentRedirection: "RadioNetwork : Redirection",
	ngapie.CauseRadioNetworkPresentResourcesNotAvailableForTheSlice: "RadioNetwork : " +
		"ResourcesNotAvailableForTheSlice",
	ngapie.CauseRadioNetworkPresentUeMaxIntegrityProtectedDataRateReason: "RadioNetwork : " +
		"UeMaxIntegrityProtectedDataRateReason",
	ngapie.CauseRadioNetworkPresentReleaseDueToCnDetectedMobility: "RadioNetwork : ReleaseDueToCnDetectedMobility",
	ngapie.CauseRadioNetworkPresentN26InterfaceNotAvailable:       "RadioNetwork : N26InterfaceNotAvailable",
	ngapie.CauseRadioNetworkPresentReleaseDueToPreEmption:         "RadioNetwork : ReleaseDueToPreEmption",

	// Values the previous ngapType package did not define. Adding them is
	// additive: they used to be reported as UNKNOWN_CAUSE_ERR.
	ngapie.CauseRadioNetworkPresentMultipleLocationReportingReferenceIDInstances: "RadioNetwork : " +
		"MultipleLocationReportingReferenceIDInstances",
	ngapie.CauseRadioNetworkPresentRsnNotAvailableForTheUp:    "RadioNetwork : RsnNotAvailableForTheUp",
	ngapie.CauseRadioNetworkPresentNpnAccessDenied:            "RadioNetwork : NpnAccessDenied",
	ngapie.CauseRadioNetworkPresentCagOnlyAccessDenied:        "RadioNetwork : CagOnlyAccessDenied",
	ngapie.CauseRadioNetworkPresentInsufficientUeCapabilities: "RadioNetwork : InsufficientUeCapabilities",
	ngapie.CauseRadioNetworkPresentRedcapUeNotSupported:       "RadioNetwork : RedcapUeNotSupported",
	ngapie.CauseRadioNetworkPresentUnknownMBSSessionID:        "RadioNetwork : UnknownMBSSessionID",
	ngapie.CauseRadioNetworkPresentIndicatedMBSSessionAreaInformationNotServedByTheGNB: "RadioNetwork : " +
		"IndicatedMBSSessionAreaInformationNotServedByTheGNB",
	ngapie.CauseRadioNetworkPresentInconsistentSliceInfoForTheSession: "RadioNetwork : InconsistentSliceInfoForTheSession",
	ngapie.CauseRadioNetworkPresentMisalignedAssociationForMulticastUnicast: "RadioNetwork : " +
		"MisalignedAssociationForMulticastUnicast",
}

var causeTransportStr = map[aper.Enumerated]string{
	ngapie.CauseTransportPresentTransportResourceUnavailable: "Transport : TransportResourceUnavailable",
	ngapie.CauseTransportPresentUnspecified:                  "Transport : Unspecified",
}

var causeNasStr = map[aper.Enumerated]string{
	ngapie.CauseNasPresentNormalRelease:         "Nas : NormalRelease",
	ngapie.CauseNasPresentAuthenticationFailure: "Nas : AuthenticationFailure",
	ngapie.CauseNasPresentDeregister:            "Nas : Deregister",
	ngapie.CauseNasPresentUnspecified:           "Nas : Unspecified",

	// Not defined by the previous ngapType package, see above.
	ngapie.CauseNasPresentUENotInPLMNServingArea: "Nas : UENotInPLMNServingArea",
}

var causeProtocolStr = map[aper.Enumerated]string{
	ngapie.CauseProtocolPresentTransferSyntaxError:                "Protocol : TransferSyntaxError",
	ngapie.CauseProtocolPresentAbstractSyntaxErrorReject:          "Protocol : AbstractSyntaxErrorReject",
	ngapie.CauseProtocolPresentAbstractSyntaxErrorIgnoreAndNotify: "Protocol : AbstractSyntaxErrorIgnoreAndNotify",
	ngapie.CauseProtocolPresentMessageNotCompatibleWithReceiverState: "Protocol : " +
		"MessageNotCompatibleWithReceiverState",
	ngapie.CauseProtocolPresentSemanticError: "Protocol : SemanticError",
	ngapie.CauseProtocolPresentAbstractSyntaxErrorFalselyConstructedMessage: "Protocol : " +
		"AbstractSyntaxErrorFalselyConstructedMessage",
	ngapie.CauseProtocolPresentUnspecified: "Protocol : Unspecified",
}

var causeMiscStr = map[aper.Enumerated]string{
	ngapie.CauseMiscPresentControlProcessingOverload:             "Misc : ControlProcessingOverload",
	ngapie.CauseMiscPresentNotEnoughUserPlaneProcessingResources: "Misc : NotEnoughUserPlaneProcessingResources",
	ngapie.CauseMiscPresentHardwareFailure:                       "Misc : HardwareFailure",
	ngapie.CauseMiscPresentOmIntervention:                        "Misc : OmIntervention",
	// Upstream renamed this value to CauseMiscPresentUnknownPLMNOrSNPN, the
	// label value stays as it was to keep existing queries working.
	ngapie.CauseMiscPresentUnknownPLMNOrSNPN: "Misc : UnknownPLMN",
	ngapie.CauseMiscPresentUnspecified:       "Misc : Unspecified",
}

func getCauseStr(causeStr map[aper.Enumerated]string, value aper.Enumerated) string {
	if str, ok := causeStr[value]; ok {
		return str
	}

	return UNKNOWN_CAUSE_ERR
}

func GetCauseErrorStr(cause *ngapie.Cause) string {
	if cause == nil {
		return UNKNOWN_NGAP_TYPE_CAUSE_ERR
	}

	// Only the pointer types implement ngapie.CauseAlt: Read has a pointer
	// receiver, so a cause group value cannot be stored in Choice.
	switch causeGroup := cause.Choice.(type) {
	case *ngapie.CauseRadioNetwork:
		return getCauseStr(causeRadioNetworkStr, causeGroup.Value)
	case *ngapie.CauseTransport:
		return getCauseStr(causeTransportStr, causeGroup.Value)
	case *ngapie.CauseNas:
		return getCauseStr(causeNasStr, causeGroup.Value)
	case *ngapie.CauseProtocol:
		return getCauseStr(causeProtocolStr, causeGroup.Value)
	case *ngapie.CauseMisc:
		return getCauseStr(causeMiscStr, causeGroup.Value)
	case *ngapie.ProtocolIESingleContainerCauseExtIEs:
		return CHOICE_EXTENSIONS_UNKNOWN_ERR
	}

	return UNKNOWN_NGAP_TYPE_CAUSE_ERR
}
