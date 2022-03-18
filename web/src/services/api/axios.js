import axios from "axios";

const axiosClient = axios.create({
  baseURL: process.env.REACT_APP_API_ENDPOINT,
});

axiosClient.interceptors.response.use(
  (response) => response.data,
  (error) => {
    console.log(error.response);
    return {
      errorCode: -1,
    };
  }
);

export default axiosClient;
