<template>
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-6 sm:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="text-xl font-semibold text-primary">OpenAI Forward</div>
        <div class="flex items-center space-x-8">
          <button v-if="!hasToken" @click.prevent="handleLogin"
                  class="bg-primary text-white px-6 py-2 !rounded-button hover:bg-blue-700 transition-colors whitespace-nowrap">
            登入
          </button>

          <button v-if="hasToken" @click.prevent="handleShowToken"
                  class="bg-primary text-white px-6 py-2 !rounded-button hover:bg-blue-700 transition-colors whitespace-nowrap">
            查看Token
          </button>
        </div>
      </div>
    </div>
  </nav>
</template>

<script>
export default {
  name: "nav",

  props: {
    hasToken: {
      type: Boolean,
      default: false
    }
  },

  setup(props, {emit}) {
    const handleLogin = () => {
      emit('showLogin');
    };

    const handleShowToken = () => {
      emit('showToken');
    };


    return {
      handleLogin,
      handleShowToken,
    }
  }
}
</script>