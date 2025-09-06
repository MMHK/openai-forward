<template>
  <div id="toast"
       class="fixed bottom-4 right-4 bg-gray-800 text-white px-4 py-2 rounded-lg transform transition-all duration-300 z-[100]"
       :class="{
         'translate-y-0 opacity-100': modelValue,
         'translate-y-full opacity-0': !modelValue
       }">
    {{ message }}
  </div>
</template>

<script>
import {watch} from "vue";

let timeoutID = 0;

export default {
  name: "toast",

  props: {
    modelValue: {
      type: Boolean,
      default: false
    },
    message: {
      type: String,
      default: ''
    }
  },

  setup(props, {emit}) {
    watch(() => props.modelValue, (value) => {
      if (value) {
        clearTimeout(timeoutID);
        timeoutID = setTimeout(() => {
          props.modelValue = false;
          emit('update:modelValue', false);
        }, 3000);
      }
    })
  }
}
</script>