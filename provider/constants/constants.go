package constants

type MessageType string

const JoinMessage MessageType = "join"
const StartMessage MessageType = "start"
const SDPMessage MessageType = "sdp"
const IceCandidateMessage MessageType = "ice-candidate"

const KeyUp = "KEYUP"
const KeyDown = "KEYDOWN"
const MouseMove = "MOUSEMOVE"
const MouseUp = "MOUSEUP"
const MouseDown = "MOUSEDOWN"
