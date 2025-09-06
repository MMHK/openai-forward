<template>
  <div id="app">
    <top-nav @show-login="handleNavigation({login: true})" @show-token="handleNavigation({token: true})"
        :has-token="hasToken"></top-nav>
    <section class="hero-section relative bg-gradient-to-br from-blue-50 to-white min-h-screen flex items-center">
      <div class="absolute inset-0 bg-white bg-opacity-70"></div>
      <div class="relative w-full max-w-7xl mx-auto px-6 sm:px-8">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-12 items-center">
          <div class="space-y-8">
            <h1 class="text-5xl font-bold text-gray-900 leading-tight">安全便捷的<br>AI 服務中轉站</h1>
            <p class="text-xl text-gray-600 leading-relaxed">三步輕鬆開始使用 AI 服務：登入、獲取 Token、連接應用</p>
            <div class="flex flex-col sm:flex-row gap-4">
              <button @click.prevent="scrollToLastSection" class="bg-primary text-white px-8 py-4 text-lg !rounded-button hover:bg-blue-700 transition-colors whitespace-nowrap">
                開始使用
              </button>
            </div>
          </div>
          <div class="flex flex-col space-y-6">
            <div class="bg-white p-6 shadow-lg border border-gray-200">
              <div class="flex items-center space-x-4 mb-4">
                <div class="w-8 h-8 bg-primary text-white flex items-center justify-center text-lg font-bold">
                  1
                </div>
                <h3 class="text-lg font-semibold">OAuth 登入</h3>
              </div>
              <p class="text-gray-600">使用您的帳號安全登入，無需記憶額外密碼</p>
            </div>
            <div class="bg-white p-6 shadow-lg border border-gray-200">
              <div class="flex items-center space-x-4 mb-4">
                <div class="w-8 h-8 bg-primary text-white flex items-center justify-center text-lg font-bold">
                  2
                </div>
                <h3 class="text-lg font-semibold">獲取臨時 Token</h3>
              </div>
              <p class="text-gray-600">選擇 OpenAI 或 Azure 端點，獲取安全的臨時憑證</p>
            </div>
            <div class="bg-white p-6 shadow-lg border border-gray-200">
              <div class="flex items-center space-x-4 mb-4">
                <div class="w-8 h-8 bg-primary text-white flex items-center justify-center text-lg font-bold">
                  3
                </div>
                <h3 class="text-lg font-semibold">連接應用</h3>
              </div>
              <p class="text-gray-600">在您喜愛的 AI 應用中使用 Token 開始創作</p>
            </div>
          </div>
        </div>
      </div>
    </section>

    <section id="apps" class="py-20 bg-gray-50">
      <div class="max-w-7xl mx-auto px-6 sm:px-8">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 mb-4">推薦應用場景</h2>
          <p class="text-xl text-gray-600">根據不同工作需求，選擇最適合的 MCP 服務組合</p>
        </div>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
          <div class="bg-gradient-to-br from-pink-400/10 to-purple-500/30 p-6 shadow-lg flex flex-wrap items-center md:justify-start md:items-start md:flex-col
            border border-gray-200/20 hover:shadow-xl transition-shadow">
            <div class="w-16 h-16 flex items-center justify-center mb-4">
              <i class="ri-palette-line ri-2x text-pink-500"></i>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3 flex-1">MCP 創意設計套件</h3>
            <div class="space-y-2 mb-4 w-[100%]">
              <div class="text-sm text-gray-600">包含服務：設計生成、風格轉換</div>
              <div class="text-sm text-gray-600">適用場景：品牌設計、UI/UX</div>
            </div>
            <p class="text-gray-600 text-sm w-[100%]">提供完整的設計工作流程，從創意發想到最終輸出</p>
          </div>
          <div class="bg-gradient-to-br from-blue-400/10 to-indigo-500/30 p-6 shadow-lg border flex flex-wrap items-center md:justify-start md:items-start md:flex-col
            border-gray-200/20 hover:shadow-xl transition-shadow">
            <div class="w-16 h-16 flex items-center justify-center mb-4">
              <i class="ri-quill-pen-line ri-2x text-blue-500"></i>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">MCP 文字創作套件</h3>
            <div class="space-y-2 mb-4">
              <div class="text-sm text-gray-600">包含服務：內容生成、多語言處理</div>
              <div class="text-sm text-gray-600">適用場景：文案、報告、翻譯</div>
            </div>
            <p class="text-gray-600 text-sm">智能文字處理解決方案，提升寫作效率與品質</p>
          </div>
          <div class="bg-gradient-to-br from-green-400/10 to-teal-500/30 p-6 shadow-lg flex flex-wrap items-center md:justify-start md:items-start md:flex-col
            border border-gray-200/20 hover:shadow-xl transition-shadow">
            <div class="w-16 h-16 flex items-center justify-center mb-4">
              <i class="ri-bar-chart-line ri-2x text-green-500"></i>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">MCP 數據分析套件</h3>
            <div class="space-y-2 mb-4">
              <div class="text-sm text-gray-600">包含服務：數據處理、智能分析</div>
              <div class="text-sm text-gray-600">適用場景：商業分析、研究報告</div>
            </div>
            <p class="text-gray-600 text-sm">深度數據分析與可視化，助您做出明智決策</p>
          </div>
          <div class="bg-gradient-to-br from-orange-400/10 to-red-500/30 p-6 shadow-lg flex flex-wrap items-center md:justify-start md:items-start md:flex-col
            border border-gray-200/20 hover:shadow-xl transition-shadow">
            <div class="w-16 h-16 flex items-center justify-center mb-4">
              <i class="ri-graduation-cap-line ri-2x text-orange-500"></i>
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-3">MCP 教育助手套件</h3>
            <div class="space-y-2 mb-4">
              <div class="text-sm text-gray-600">包含服務：智能輔導、知識管理</div>
              <div class="text-sm text-gray-600">適用場景：在線學習、知識庫</div>
            </div>
            <p class="text-gray-600 text-sm">打造個性化學習體驗，提升學習效果</p>
          </div>
        </div>
      </div>
    </section>
    <section class="py-20 bg-white">
      <div class="max-w-4xl mx-auto px-6 sm:px-8">
        <div class="text-center mb-16">
          <h2 class="text-4xl font-bold text-gray-900 mb-4">常見問題解答</h2>
          <p class="text-xl text-gray-600">快速了解 OpenAI Forward 的使用方式</p>
        </div>
        <div class="space-y-6">
          <div class="border-gray-200 p-6 shadow-lg">
            <h3 class="text-lg font-semibold text-gray-900 mb-3">如何選擇合適的服務？</h3>
            <p class="text-gray-600">根據您的具體需求選擇：文字處理選擇 ChatGPT 或 Claude，圖像生成選擇 DALL-E 或 Midjourney，語音處理選擇
              Whisper。</p>
          </div>
          <div class="border-gray-200 p-6 shadow-lg">
            <h3 class="text-lg font-semibold text-gray-900 mb-3">Token 有效期多久？</h3>
            <p class="text-gray-600">臨時 Token 的有效期為 24 小時，到期後需要重新登入獲取新的 Token。這樣設計是為了確保您的帳戶安全。</p>
          </div>
          <div class="border-gray-200 p-6 shadow-lg">
            <h3 class="text-lg font-semibold text-gray-900 mb-3">如何續期 Token？</h3>
            <p class="text-gray-600">當 Token 即將到期時，系統會自動提醒您。您只需要重新登入即可獲取新的 Token，整個過程非常簡單。</p>
          </div>
          <div class="border-gray-200 p-6 shadow-lg">
            <h3 class="text-lg font-semibold text-gray-900 mb-3">支援哪些客戶端？</h3>
            <p class="text-gray-600">支援所有相容 OpenAI API 的客戶端，包括 Cherry Studio、Notion AI、各種瀏覽器擴展和移動應用程式。</p>
          </div>
        </div>
      </div>
    </section>
    <section id="cta_section" class="py-20 bg-primary text-white">
      <div class="max-w-4xl mx-auto px-6 sm:px-8 text-center">
        <h2 class="text-4xl font-bold mb-6">立即開始使用 AI 服務</h2>
        <p class="text-xl mb-8 opacity-90">安全、便捷、高效的 AI 服務中轉解決方案</p>
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <button v-if="!hasToken"
                  @click="handleNavigation({login:true})"
                  class="bg-white text-primary px-8 py-4 text-lg !rounded-button hover:bg-gray-100 transition-colors whitespace-nowrap font-semibold">
            立即登入
          </button>
          <button v-if="hasToken"
                  @click="handleNavigation({token:true})"
              class="bg-white text-primary px-8 py-4 text-lg !rounded-button hover:bg-gray-100 transition-colors whitespace-nowrap font-semibold">
            查看 Token
          </button>
          <a target="_blank" href="https://github.com/MMHK/openai-forward" class="border-2 border-white text-white px-8 py-4 text-lg !rounded-button hover:bg-white hover:text-primary transition-colors whitespace-nowrap">
            Github
          </a>
        </div>
      </div>
    </section>
    <footer class="bg-gray-900 text-white py-5">
      <div class="max-w-7xl mx-auto px-6 sm:px-8">
        <div class="text-center text-sm text-gray-400 py-4">
          <p>&copy; {{ new Date().getFullYear() }} MixMedia Production Limited. 保留所有權利。</p>
        </div>
      </div>
    </footer>

    <login-popup
        @close-popup="handleNavigation({login:false})"
        v-if="showLoginPopup"></login-popup>
    <token-popup
        @close-popup="handleNavigation({token:false})"
        v-if="showTokenPopup"></token-popup>
    <toast :message="ToastInfo.message" v-model="ToastInfo.show"></toast>
  </div>
</template>

<script>
import TopNav from "./components/nav.vue";
import Toast from "./components/toast.vue";
import LoginPopup from "./components/login-popup.vue";
import TokenPopup from "./components/token-popup.vue";
import {computed, provide, ref} from "vue";

const SHOW_TOAST = 'showToast';

export default {
  name: "App",

  components: {
    TopNav,
    Toast,
    LoginPopup,
    TokenPopup
  },

  setup() {
    const ToastInfo = ref({
      show: false,
      message: ''
    });

    const showToast = (message) => {
      ToastInfo.value.show = true;
      ToastInfo.value.message = message;
      setTimeout(() => {
        ToastInfo.value.show = false;
        ToastInfo.value.message = '';
      }, 3000);
    }
    provide(SHOW_TOAST, showToast);



    const hasToken = computed(() => {
      return !!(localStorage.getItem('token'));
      // return true;
    })

    const showLoginPopup = ref(false);
    const showTokenPopup = ref(false);

    const handleNavigation = (options={login: false, token: false}) => {
      console.log(options);

      showTokenPopup.value = options.token;
      showLoginPopup.value = options.login;
    }

    const scrollToLastSection = () => {
      const lastSection = document.getElementById('cta_section');
      console.log(lastSection);
      if (lastSection) {
        lastSection.scrollIntoView({behavior: 'smooth'});
      }
    }

    return {
      ToastInfo,
      showToast,
      hasToken,
      showLoginPopup,
      showTokenPopup,
      handleNavigation,
      scrollToLastSection,
    }
  }
}
</script>