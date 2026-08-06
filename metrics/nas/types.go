package nas

const (
	SUBSYSTEM_NAME = "nas"

	NAS_MSG_RCV_COUNTER_NAME = "msg_received_total"
	NAS_MSG_RCV_COUNTER_DESC = "Total number of received NAS message"

	NAS_MSG_SENT_COUNTER_NAME = "nas_msg_sent_total"
	NAS_MSG_SENT_COUNTER_DESC = "Total number of NAS message sent"
)

// metric collectors label values
const (
	NAME_LABEL   = "name"
	STATUS_LABEL = "status"
	CAUSE_LABEL  = "cause"
)

// These values are tied to the metrics NAS message type.
//
// They are the public label vocabulary of the NAS metrics: renaming any of
// them breaks existing dashboards and alerting rules, so they are pinned here
// rather than derived from the message.MsgType names of github.com/free5gc/nas.
// msgTypeToMetricName in message.go maps upstream message types onto them.
const (
	UNKNOWN_GMM_MESSAGE                                  = "Unknown gmm message"
	AUTHENTICATION_REQUEST                               = "AuthenticationRequest"
	AUTHENTICATION_RESPONSE                              = "AuthenticationResponse"
	AUTHENTICATION_RESULT                                = "AuthenticationResult"
	AUTHENTICATION_FAILURE                               = "AuthenticationFailure"
	AUTHENTICATION_REJECT                                = "AuthenticationReject"
	REGISTRATION_REQUEST                                 = "RegistrationRequest"
	REGISTRATION_ACCEPT                                  = "RegistrationAccept"
	REGISTRATION_ACCEPT_TIMER                            = "RegistrationAcceptTimer"
	REGISTRATION_COMPLETE                                = "RegistrationComplete"
	REGISTRATION_REJECT                                  = "RegistrationReject"
	UL_NAS_TRANSPORT                                     = "ULNASTransport"
	DL_NAS_TRANSPORT                                     = "DLNASTransport"
	DEREGISTRATION_REQUEST_UE_ORIGINATING_DEREGISTRATION = "DeregistrationRequestUEOriginatingDeregistration"
	DEREGISTRATION_ACCEPT_UE_ORIGINATING_DEREGISTRATION  = "DeregistrationAcceptUEOriginatingDeregistration"
	DEREGISTRATION_REQUEST_UE_TERMINATED_DEREGISTRATION  = "DeregistrationRequestUETerminatedDeregistration"
	DEREGISTRATION_ACCEPT_UE_TERMINATED_DEREGISTRATION   = "DeregistrationAcceptUETerminatedDeregistration"
	SERVICE_REQUEST                                      = "ServiceRequest"
	SERVICE_ACCEPT                                       = "ServiceAccept"
	SERVICE_REJECT                                       = "ServiceReject"
	CONFIGURATION_UPDATE_COMMAND                         = "ConfigurationUpdateCommand"
	CONFIGURATION_UPDATE_COMMAND_TIMER                   = "ConfigurationUpdateCommandTimer"
	CONFIGURATION_UPDATE_COMPLETE                        = "ConfigurationUpdateComplete"
	IDENTITY_REQUEST                                     = "IdentityRequest"
	IDENTITY_RESPONSE                                    = "IdentityResponse"
	NOTIFICATION                                         = "Notification"
	NOTIFICATION_TIMER                                   = "NotificationTimer"
	NOTIFICATION_RESPONSE                                = "NotificationResponse"
	SECURITY_MODE_COMMAND                                = "SecurityModeCommand"
	SECURITY_MODE_COMPLETE                               = "SecurityModeComplete"
	SECURITY_MODE_REJECT                                 = "SecurityModeReject"
	SECURITY_PROTECTED_5GS_NAS_MESSAGE                   = "SecurityProtected5GSNASMessage"
	STATUS_5GMM                                          = "Status5GMM"

	// 5GMM message types that the previous nas API did not expose and that
	// therefore used to be counted as UNKNOWN_GMM_MESSAGE. Adding them is
	// additive: no pre-existing label value changes meaning.
	CONTROL_PLANE_SERVICE_REQUEST                  = "ControlPlaneServiceRequest"
	NETWORK_SLICE_SPECIFIC_AUTHENTICATION_COMMAND  = "NetworkSliceSpecificAuthenticationCommand"
	NETWORK_SLICE_SPECIFIC_AUTHENTICATION_COMPLETE = "NetworkSliceSpecificAuthenticationComplete"
	NETWORK_SLICE_SPECIFIC_AUTHENTICATION_RESULT   = "NetworkSliceSpecificAuthenticationResult"
	RELAY_KEY_REQUEST                              = "RelayKeyRequest"
	RELAY_KEY_ACCEPT                               = "RelayKeyAccept"
	RELAY_KEY_REJECT                               = "RelayKeyReject"
	RELAY_AUTHENTICATION_REQUEST                   = "RelayAuthenticationRequest"
	RELAY_AUTHENTICATION_RESPONSE                  = "RelayAuthenticationResponse"
)

// Additional error causes
const (
	RAN_UE_NIL_ERR      = "RanUe is nil"
	AMF_UE_NIL_ERR      = "AmfUe is nil"
	RAN_NIL_ERR         = "Ran is nil"
	NAS_PDU_NIL_ERR     = "nasPdu is nil"
	AUTH_CTX_UE_NIL_ERR = "Authentication Context of UE is nil"
	NAS_MSG_BUILD_ERR   = "Could not build NAS message"
	DECODE_NAS_MSG_ERR  = "Could not decode NAS message"
	AUSF_AUTH_ERR       = "Ausf Authentication Failure"
	HRES_AUTH_ERR       = "HRES* validation failure"
)

var nasMetricsEnabled bool

func IsNasMetricsEnabled() bool {
	return nasMetricsEnabled
}

func EnableNasMetrics() {
	nasMetricsEnabled = true
}
