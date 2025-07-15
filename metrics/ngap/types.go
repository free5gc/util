package ngap

const (
	SUBSYSTEM_NAME = "ngap"

	MSG_RCV_COUNTER_NAME = "msg_received_total"
	MSG_RCV_COUNTER_DESC = "Total number of received NGAP message"

	MSG_SENT_COUNTER_NAME = "msg_sent_total"
	MSG_SENT_COUNTER_DESC = "Total number of NGAP message sent"
)

// metric collectors label values
const (
	NAME_LABEL   = "name"
	STATUS_LABEL = "status"
	CAUSE_LABEL  = "cause"
)

const (
	AMF_CONFIGURATION_UPDATE                   = "AMFConfigurationUpdate"
	AMF_CONFIGURATION_UPDATE_ACKNOWLEDGE       = "AMFConfigurationUpdateAcknowledge"
	AMF_CONFIGURATION_UPDATE_FAILURE           = "AMFConfigurationUpdateFailure"
	AMF_STATUS_INDICATION                      = "AMFStatusIndication"
	DEACTIVATE_TRACE                           = "DeactivateTrace"
	DOWNLINK_NAS_TRANSPORT                     = "DownlinkNasTransport"
	DOWNLINK_NON_UE_ASSOCIATED_NRPPA_TRANSPORT = "DownlinkNonUEAssociatedNRPPATransport"
	DOWNLINK_RAN_CONFIGURATION_TRANSFER        = "DownlinkRanConfigurationTransfer"
	DOWNLINK_RAN_STATUS_TRANSFER               = "DownlinkRanStatusTransfer"
	DOWNLINK_UE_ASSOCIATED_NRPPA_TRANSPORT     = "DownlinkUEAssociatedNRPPaTransport"
	ERROR_INDICATION                           = "ErrorIndication"
	ERROR_INDICATION_WITH_SCTP_CONN            = "ErrorIndicationWithSctpConn"
	HANDOVER_CANCEL_ACKNOWLEDGE                = "HandoverCancelAcknowledge"
	HANDOVER_COMMAND                           = "HandoverCommand"
	HANDOVER_PREPARATION_FAILURE               = "HandoverPreparationFailure"
	HANDOVER_REQUEST                           = "HandoverRequest"
	INITIAL_CONTEXT_SETUP_REQUEST              = "InitialContextSetupRequest"
	INITIAL_CONTEXT_SETUP_RESPONSE             = "InitialContextSetupResponse"
	INITIAL_UE_MESSAGE                         = "InitialUeMessage"
	LOCATION_REPORTING_CONTROL                 = "LocationReportingControl"
	NAS_NON_DELIVERY_INDICATION                = "NasNonDeliveryIndication"
	NG_RESET                                   = "NGReset"
	NG_RESET_ACKNOWLEDGE                       = "NGResetAcknowledge"
	NG_SETUP_FAILURE                           = "NGSetupFailure"
	NG_SETUP_REQUEST                           = "NGSetupRequest"
	NG_SETUP_RESPONSE                          = "NGSetupResponse"
	OVERLOAD_START                             = "OverloadStart"
	OVERLOAD_STOP                              = "OverloadStop"
	PAGING                                     = "Paging"
	PATH_SWITCH_REQUEST_ACKNOWLEDGE            = "PathSwitchRequestAcknowledge"
	PATH_SWITCH_REQUEST_FAILURE                = "PathSwitchRequestFailure"
	PDUSESSION_RESOURCE_MODIFY_CONFIRM         = "PDUSessionResourceModifyConfirm"
	PDUSESSION_RESOURCE_MODIFY_INDICATION      = "PDUSessionResourceModifyIndication"
	PDUSESSION_RESOURCE_MODIFY_REQUEST         = "PDUSessionResourceModifyRequest"
	PDUSESSION_RESOURCE_MODIFY_RESPONSE        = "PDUSessionResourceModifyResponse"
	PDUSESSION_RESOURCE_NOTIFY                 = "PDUSessionResourceNotify"
	PDUSESSION_RESOURCE_RELEASE_COMMAND        = "PDUSessionResourceReleaseCommand"
	PDUSESSION_RESOURCE_RELEASE_RESPONSE       = "PDUSessionResourceReleaseResponse"
	PDUSESSION_RESOURCE_SETUP_REQUEST          = "PDUSessionResourceSetupRequest"
	PDUSESSION_RESOURCE_SETUP_RESPONSE         = "PDUSessionResourceSetupResponse"
	RAN_CONFIGURATION_UPDATE_ACKNOWLEDGE       = "RanConfigurationUpdateAcknowledge"
	RAN_CONFIGURATION_UPDATE_FAILURE           = "RanConfigurationUpdateFailure"
	RAN_CONFIGURATION_UPDATE_UPDATE            = "RanConfigurationUpdateUpdate"
	REROUTE_NAS_REQUEST                        = "RerouteNasRequest"
	UE_CONTEXT_MODIFICATION_FAILURE            = "UEContextModificationFailure"
	UE_CONTEXT_MODIFICATION_REQUEST            = "UEContextModificationRequest"
	UE_CONTEXT_MODIFICATION_RESPONSE           = "UEContextModificationResponse"
	UE_CONTEXT_RELEASE_COMMAND                 = "UEContextReleaseCommand"
	UE_CONTEXT_RELEASE_COMPLETE                = "UEContextReleaseComplete"
	UE_CONTEXT_RELEASE_REQUEST                 = "UEContextReleaseRequest"
	UE_RADIO_CAPABILITY_CHECK_REQUEST          = "UERadioCapabilityCheckRequest"
	UE_RADIO_CAPABILITY_CHECK_RESPONSE         = "UERadioCapabilityCheckResponse"
	UE_TNLA_BINDING_RELEASE_REQUEST            = "UETNLABindingReleaseRequest"
	UPLINK_NAS_TRANSPORT                       = "UplinkNasTransport"
)

// Additional error causes
const (
	AMF_TRAFFIC_LOAD_REDUCTION_INDICATION_OOO_ERR = "AmfTrafficLoadReductionIndication out of range (should be 1 ~ 99)"
	AMF_TIME_REINIT_ERR                           = "Please Wait at least for the indicated time before reinitiating " +
		"toward same AMF"
	AMF_UE_NIL_ERR                                 = "AmfUe is nil"
	AOI_LIST_OOR_ERR                               = "AOI List out of range"
	CAUSE_NIL_ERR                                  = "Cause present is nil"
	ERROR_INDICATION_CAUSE_AND_CRITICALITY_NIL_ERR = "Both cause and criticality are nil, one of them should be set " +
		"in the message"
	GUAMI_LIST_OOR_ERR                      = "GUAMI List out of range"
	HANDOVER_REQUIRED_DUP_ERR               = "Handover Required Duplicated"
	LOCATION_REPORTING_REFERENCE_ID_OOR_ERR = "LocationReportingReferenceIDToBeCancelled out of range " +
		"(should be 1 ~ 64)"
	NAS_PDU_NIL_ERR                             = "Nas Pdu is nil"
	NGAP_MSG_BUILD_ERR                          = "Could not build NAS message"
	NGAP_MSG_NIL_ERR                            = "Ngap Message is nil"
	NF_INTERFACE_LEN_ZERO_ERR                   = "length of partOfNGInterface is 0"
	NRPPA_LEN_ZERO_ERR                          = "length of NRPPA-PDU is 0"
	NSSAI_LIST_OOR_ERR                          = "NSSAI List out of range"
	PDU_SESS_RESOURCE_MODIFY_LIST_NIL_ERR       = "PDU Session Resource Modify List indication is nil"
	PDU_SESS_RESOURCE_RELEASED_LIST_NIL_ERR     = "Pdu Session Resource Release List is nil"
	PDU_SESS_RESOURCE_SWITCH_OOO_ERR            = "Pdu Session Resource Switched List out of range"
	PDU_LIST_OOR_ERR                            = "Pdu List out of range"
	RAN_NIL_ERR                                 = "Ran is nil"
	RAN_UE_NIL_ERR                              = "RanUe is nil"
	SCTP_SOCKET_WRITE_ERR                       = "Write to SCTP socket failed"
	SOURCE_UE_NIL_ERR                           = "SourceUe is nil"
	SRC_TO_TARGET_TRANSPARENT_CONTAINER_NIL_ERR = "Source To Target TransparentContainer is nil"
	TARGET_RAN_NIL_ERR                          = "targetRan is nil"
	UE_CTX_NIL                                  = "Ue context is nil"
	UNKNOWN_NGAP_TYPE_CAUSE_ERR                 = "unknown ngapType.Cause"
)

var ngapMetricsEnabled bool

func IsNgapMetricsEnabled() bool {
	return ngapMetricsEnabled
}

func EnableNgapMetrics() {
	ngapMetricsEnabled = true
}
