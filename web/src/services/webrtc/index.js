export const addRemoteSdp = async (pc, offer) => {
  await pc.setRemoteDescription(new RTCSessionDescription(offer));
  const answer = await pc.createAnswer();
  await pc.setLocalDescription(answer);
  return answer;
};

export const addIceCandidate = async (pc, iceCandidate) => {
  try {
    pc.addIceCandidate(new RTCIceCandidate(iceCandidate));
  } catch (e) {
    console.log(e);
  }
};
