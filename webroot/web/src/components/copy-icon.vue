<template>
  <span @click.prevent="copyText">
    <i v-if="!copied" class="ri-file-copy-line"></i>
    <i v-if="copied" class="ri-check-line"></i>
  </span>
</template>

<script>
import {inject, ref} from "vue";

export default {
  name: "copy-icon",

  props: {
    text: {
      type: String,
      default: ''
    }
  },

  setup() {
    const copied = ref(false);
    const showToast = inject('showToast');

    const copyText = (text) => {
      return navigator.clipboard.writeText(text).then(() => {
        copied.value = true;
        showToast('端點地址已複製到剪貼簿');
        setTimeout(() => {
          copied.value = false;
        }, 2000);
      })
    }

    return {
      copied,
      copyText,
    }
  }
}
</script>