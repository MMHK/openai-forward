import axios from "axios";

const BASE_URL = global.API_ENDPOINT || "/api/v1";
const http = axios.create({
  baseURL: BASE_URL,
});

class HttpError extends Error {
  constructor(message, status) {
    super(message);
    this.status = status;
  }
}
class AuthError extends HttpError {}

export {HttpError,AuthError}

http.interceptors.request.use((config) => {
    const apikey = localStorage.getItem("token");
    if (apikey) {
        config.headers.Authorization = `Bearer ${apikey}`;
    }
    return config;
})

http.interceptors.response.use((response) => {
    const json = response.data;
    if (response.status == 401) {
        return Promise.reject(new AuthError("Unauthorized", response.status));
    }
    if (json.status && json.data) {
      return Promise.resolve(json.data);
    }
    if (json.error) {
        return Promise.reject(new HttpError(json.error, response.status));
    }
    return Promise.reject(new HttpError("Unknown error", response.status));
}, (error) => {
    if (error.response.status == 401) {
        return Promise.reject(new AuthError("Unauthorized", error.response.status));
    }
    return Promise.reject(error);
})

export default {
    async auth() {
        let redirect = window.location.href;
        window.location.href = `${BASE_URL}/auth?redirect=${encodeURI(redirect)}`;
    },
    async authCallback(code) {
        const apikey = await http.get(`/auth/callback?code=${code}"`)
        const {key, expire_at} = apikey;
        localStorage.setItem("token", key);
        localStorage.setItem("token_expires_at", expire_at);
        window.location.href = "/";
    },

    async AzureModels() {
        return http.get("/azure/models");
    },
    async OpenAIModel() {
        return http.get(`/openai/models`);
    },

    GetTokenInfo() {
        return {
            token: localStorage.getItem("token"),
            token_expires_at: localStorage.getItem("token_expires_at")
        }
    }
}
