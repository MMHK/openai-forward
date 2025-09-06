import "remixicon/fonts/remixicon.css";
import "./style/tailwindcss.css";
import "./style/main.scss"
import httpService from "./service";
import {createApp} from "vue";
import App from "./App.vue";

createApp(App).mount("#app");

const query = new URLSearchParams(window.location.search);
if (query.has('code')) {
    const code = query.get('code');
    httpService.authCallback(code);
}


// const copyToken = document.getElementById('copyToken');
// if (copyToken) {
//     copyToken.addEventListener('click', function () {
//         const tokenText = 'sk-proj-abc123def456ghi789jkl012mno345pqr678stu901vwx234yz...';
//         navigator.clipboard.writeText(tokenText).then(function () {
//             copyToken.textContent = '已複製';
//             showToast('Token 已複製到剪貼簿');
//             setTimeout(function () {
//                 copyToken.textContent = '複製';
//             }, 2000);
//         });
//     });
// }
// const modelItems = document.querySelectorAll('.flex.items-center.gap-1.px-3.py-1.bg-blue-100');
// modelItems.forEach(item => {
//     item.addEventListener('click', function () {
//         const modelName = this.querySelector('span').textContent;
//         navigator.clipboard.writeText(modelName).then(() => {
//             const icon = this.querySelector('i');
//             icon.classList.remove('ri-file-copy-line');
//             icon.classList.add('ri-check-line');
//             showToast(`${modelName} 已複製到剪貼簿`);
//             setTimeout(() => {
//                 icon.classList.remove('ri-check-line');
//                 icon.classList.add('ri-file-copy-line');
//             }, 2000);
//         });
//     });
// });
// const endpointDiv = document.querySelector('.flex.items-center.gap-1.text-sm.text-gray-600.mt-1');
// if (endpointDiv) {
//     endpointDiv.addEventListener('click', function () {
//         const endpointText = this.querySelector('span').textContent;
//         navigator.clipboard.writeText(endpointText).then(() => {
//             const icon = this.querySelector('i');
//             icon.classList.remove('ri-file-copy-line');
//             icon.classList.add('ri-check-line');
//             showToast('端點地址已複製到剪貼簿');
//             setTimeout(() => {
//                 icon.classList.remove('ri-check-line');
//                 icon.classList.add('ri-file-copy-line');
//             }, 2000);
//         });
//     });
// }
//
//
// document.addEventListener('DOMContentLoaded', function () {
//     const openaiBtn = document.getElementById('openaiBtn');
//     const azureBtn = document.getElementById('azureBtn');
//     openaiBtn.addEventListener('click', function () {
//         openaiBtn.classList.add('bg-primary', 'text-white');
//         openaiBtn.classList.remove('border', 'border-gray-300', 'text-gray-700');
//         azureBtn.classList.remove('bg-primary', 'text-white');
//         azureBtn.classList.add('border', 'border-gray-300', 'text-gray-700');
//     });
//     azureBtn.addEventListener('click', function () {
//         azureBtn.classList.add('bg-primary', 'text-white');
//         azureBtn.classList.remove('border', 'border-gray-300', 'text-gray-700');
//         openaiBtn.classList.remove('bg-primary', 'text-white');
//         openaiBtn.classList.add('border', 'border-gray-300', 'text-gray-700');
//     });
// });
