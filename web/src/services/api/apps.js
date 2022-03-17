import axiosClient from "./axios";

const APP_LIST_API = "apps";

export const getAppList = async () => {
  const device = getDevice();
  const url = `${APP_LIST_API}?device=${device}`;
  const resp = await axiosClient.get(url);
  return resp;
};

export const getDevice = () => {
  const ua = navigator.userAgent;
  if (/(tablet|ipad|playbook|silk)|(android(?!.*mobi))/i.test(ua)) {
    return "mobile"; // tablet
  } else if (
    /Mobile|Android|iP(hone|od)|IEMobile|BlackBerry|Kindle|Silk-Accelerated|(hpw|web)OS|Opera M(obi|ini)/.test(
      ua
    )
  ) {
    return "mobile";
  }
  return "pc";
};
