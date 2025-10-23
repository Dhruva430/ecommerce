import axios from "axios";

const url = "http://localhost:8080/api/";

export const apiClient = axios.create({
  baseURL: url.toString(),
  adapter: "fetch",
});
