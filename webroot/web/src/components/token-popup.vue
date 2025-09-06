<template>
  <div id="tokenPanel" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white p-8 max-w-2xl w-full mx-4 shadow-2xl">
      <div class="flex justify-between items-center mb-6">
        <h2 class="text-2xl font-bold text-gray-900">Token 管理面板</h2>
        <button @click.prevent="closePopup" class="w-8 h-8 flex items-center justify-center text-gray-500 hover:text-gray-700">
          <i class="ri-close-line ri-lg"></i>
        </button>
      </div>
      <div class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">選擇端點</label>
          <div class="flex space-x-4">
            <button @click.prevent="setActive({openai: true})"
                    :class="{'bg-primary text-white': activeInfo.openai, 'border border-gray-300 text-gray-700': activeInfo.azure}"
                    class="px-4 py-2 !rounded-button hover:bg-primary">OpenAI</button>
            <button @click.prevent="setActive({azure: true})"
                    :class="{'bg-primary text-white': activeInfo.azure, 'border border-gray-300 text-gray-700': activeInfo.openai}"
                    class="px-4 py-2 !rounded-button hover:bg-primary">
              Azure
            </button>
          </div>
        </div>
        <div class="bg-gray-50 p-4">
          <div class="flex justify-between items-center mb-2">
            <span class="text-sm font-medium text-gray-700">Token</span>
            <copy-icon :text="EndPointInfo.token" />
          </div>
          <div class="bg-white p-3 border border-gray-200 font-mono text-sm text-gray-800 break-all">
            {{ EndPointInfo.token }}
          </div>
        </div>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <span class="text-sm font-medium text-gray-700">端點</span>
            <div class="flex items-center gap-1 text-sm text-gray-600 mt-1 cursor-pointer hover:text-gray-800 group">
              <span>{{ EndPointInfo.endpoint }}</span>
              <copy-icon :text="EndPointInfo.endpoint" />
            </div>
          </div>
          <div>
            <span class="text-sm font-medium text-gray-700">到期時間</span>
            <div class="text-sm text-gray-600 mt-1">{{ EndPointInfo.expires }}</div>
          </div>
        </div>
        <div>
          <span class="text-sm font-medium text-gray-700 block mb-2">可用模型</span>
          <div class="flex flex-wrap gap-2">
            <div v-for="(model, i) in EndPointInfo.model" :key="`model-${i}`"
                class="flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-800 text-sm cursor-pointer hover:bg-blue-200 transition-colors group">
              <span>{{ model }}</span>
              <copy-icon text="model" />
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import {onMounted, ref} from "vue";
import CopyIcon from "./copy-icon.vue";
import httpService from "../service";

export default {
  name: "token-popup",

  components: {
    CopyIcon
  },

  setup(props, {emit}) {
    const closePopup = () => {
      emit('closePopup');
    };

    const EndPointInfo = ref({
      endpoint: "https://api.openai.com/v1",
      expires: "2023-09-01 12:00:00",
      token: "sk-proj-abc123def456ghi789jkl012mno345pqr678stu901vwx234yz...",
      model: [
          "GPT-4",
          "GPT-3.5-turbo",
          "DALL-E-3",
          "Whisper"
      ]
    });

    const ModelsInfo = ref({
      openai: [],
      azure: [],
    });

    onMounted(async () => {
      const [azure, openai] = await Promise.all([
          httpService.AzureModels(),
          httpService.OpenAIModel()
      ])
      ModelsInfo.value = {
        openai,
        azure
      }
    })

    const activeInfo = ref({
      openai: false,
      azure: false
    });

    const setActive = (options = {openai: false, azure: false}) => {
      activeInfo.value.azure = options.azure;
      activeInfo.value.openai = options.openai;

      const tokenInfo = httpService.GetTokenInfo();


      if (options.openai) {
        EndPointInfo.value = {
          ...EndPointInfo.value,
          endpoint: `${location.protocol}//${location.hostname}/openai`,
          token: tokenInfo.token,
          expires: tokenInfo.token_expires_at,
          model: ModelsInfo.value.openai
        }
      }

      if (options.azure) {
        EndPointInfo.value = {
          ...EndPointInfo.value,
          endpoint: `${location.protocol}//${location.hostname}/azure`,
          token: tokenInfo.token,
          expires: tokenInfo.token_expires_at,
          model: ModelsInfo.value.azure
        }
      }
    };

    setActive({openai: true})

    return {
      closePopup,
      EndPointInfo,
      activeInfo,
      setActive
    }
  }
}
</script>