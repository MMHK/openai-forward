<template>
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-[60]">
    <div class="bg-white p-8 rounded-lg max-w-md w-full mx-4 shadow-2xl">
      <div class="flex justify-between items-center mb-8">
        <h2 class="text-2xl font-bold text-gray-900">登入帳號</h2>
        <button @click.prevent="closePopup" class="w-8 h-8 flex items-center justify-center text-gray-500 hover:text-gray-700">
          <i class="ri-close-line ri-lg"></i>
        </button>
      </div>
      <button @click="handleLogin"
              class="w-full bg-[#E42527] text-white py-3 !rounded-button flex items-center justify-center gap-2 hover:bg-[#C41E20] transition-colors">
        <i class="ri-login-box-line"></i>
        <span>Login with Zoho Account</span>
      </button>
    </div>
  </div>
</template>

<script>
import HttpService from "../service";

export default {
  name: "login-popup",

  setup(props, {emit}) {
    const closePopup = () => {
      emit("closePopup");
    };

    const handleLogin = () => {
      HttpService.auth();
    }

    return {
      closePopup,
      handleLogin,
    }
  }
}
</script>