import axiosClient from "./axios";

const PROVIDER_LIST_API = "providers";

export const getProviderList = async () => {
  const resp = await axiosClient.get(PROVIDER_LIST_API);
  return resp;
};
