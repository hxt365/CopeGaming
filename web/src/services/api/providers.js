import axiosClient from "./axios";

const PROVIDER_LIST_API = "providers";

export const getProviderList = async (ownerID) => {
  let url = `${PROVIDER_LIST_API}`;
  if (ownerID !== undefined) {
    url += `?owner=${ownerID}`;
  }

  const resp = await axiosClient.get(url);
  return resp;
};
