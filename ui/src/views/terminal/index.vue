<template>
  <div class="xterm-container" :style="{
    backgroundColor: themes[settingsStore.theme.content].xterm.background,
  }">
    <div id="terminal" class="terminal" ref="terminalRef" :style="{
      backgroundColor: themes[settingsStore.theme.content].xterm.background,
    }"></div>
  </div>
</template>

<script setup>
import { Terminal } from '@xterm/xterm';
import { FitAddon } from '@xterm/addon-fit';
import { WebLinksAddon } from '@xterm/addon-web-links';
import '@xterm/xterm/css/xterm.css';
import { onBeforeMount, onMounted, ref } from 'vue';
import useSettingsStore from '../../store/modules/settings';
import useTagsView from '../../store/modules/tagsView';
import themeJson from '@/assets/theme.json'
import { useRouter, useRoute } from 'vue-router'
const terminalRef = ref(null);
const settingsStore = useSettingsStore()
const tagsView = useTagsView()
const themes = ref(themeJson);
const router = useRouter()
const route = useRoute()
//变量
let term = null
let ws = null
const termSize = {
  cols: 80,
  rows: 80
}

// 获取 URL 参数的函数
const getTerminalId = () => {
  return route.params.id; // 获取 /terminal/:id 中的 id 参数
};
//初始化终端
function initTerminal() {
  term = new Terminal({
    cursorBlink: true,
    fontSize: 12, // 减小字体大小，测试效果
    fontFamily: 'monospace', // 明确指定等宽字体
    theme: {
      background: themes.value[settingsStore.theme.content].xterm.background,
      foreground: themes.value[settingsStore.theme.content].xterm.foreground,
      cursor: themes.value[settingsStore.theme.content].xterm.cursor,
    },
  });
  //添加插件
  const fitAddon = new FitAddon();
  term.loadAddon(fitAddon);
  term.loadAddon(new WebLinksAddon());
  term.open(terminalRef.value);
  fitAddon.fit();
  termSize.cols = term.cols
  termSize.rows = term.rows
  // 监听窗口大小变化
  window.addEventListener('resize', () => {
    fitAddon.fit();
    termSize.cols = term.cols
    termSize.rows = term.rows
    sendTerminalSize()
  });
  //监听用户输入
  term.onData((data) => {
    if (ws && ws.readyState === WebSocket.OPEN) {
      ws.send(JSON.stringify({ type: "command", data: data }));
    }
  });
}
// 更新终端主题
function updateTerminalTheme() {
  if (term) {
    console.log(term);
    term.options.theme = {
      background: themes.value[settingsStore.theme.content].xterm.background,
      foreground: themes.value[settingsStore.theme.content].xterm.foreground,
      cursor: themes.value[settingsStore.theme.content].xterm.cursor,
    };
  }
}
//连接SSH
function connectSSH() {
  const terminalId = getTerminalId(); // 获取参数
  const sshId = tagsView.cachedViewMap[terminalId]
  // 使用参数构造 WebSocket URL
  ws = new WebSocket(`${import.meta.env.VITE_APP_BASE_API}/ws/${sshId}`);
  ws.onopen = () => {
    console.log("WebSocket 连接已建立");
    // 发送 SSH 连接信息
    const sshConfig = {
      type: "connect",
      id: '',
      data: JSON.stringify(termSize)
    };
    ws.send(JSON.stringify(sshConfig));
  };
  ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    if (msg.type === "info") {
      switch (msg.data) {
        case "connected":
          //连接上之后发送终端大小
          sendTerminalSize()
          break;
        default:
          break;
      }
    } else if (msg.type === "output") {
      term.write(msg.data);
    } else if (msg.type === "error") {
      term.write(`错误: ${msg.data}\r\n`);
    }
  };
  ws.onclose = () => {
    console.log("WebSocket 连接已关闭");
    term.write("\r\n连接已关闭\r\n");
  };
  ws.onerror = (error) => {
    console.error("WebSocket 错误:", error);
    term.write("\r\nWebSocket 错误\r\n");
  };
}

//发送终端行列信息
function sendTerminalSize() {
  if (ws && ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify({
      type: "resize",
      data: JSON.stringify(termSize)
    }));
  }
};


//获取ssh连接的信息
onBeforeMount(() => {
  if (tagsView.cachedViews.length === 0) {
    router.push({ path: '/' });
  }
  // 可以在这里验证参数
  const terminalId = getTerminalId();
  if (!terminalId) {
    router.push({ path: '/' }); // 如果没有参数，重定向到首页
  }
});
onMounted(async () => {
  await nextTick(); // 等待 DOM 更新
  initTerminal();
  connectSSH();
});

// 监听主题变化
watch(() => settingsStore.theme.content, (newTheme) => {
  updateTerminalTheme();
});

// 当组件重新激活时更新主题（适用于 keep-alive）
onActivated(() => {
  updateTerminalTheme();
});
</script>

<style lang="scss" scoped>
.xterm-container {
  width: 100%;
  height: calc(100vh - 50px);
  padding: 0;
  /* 移除内边距 */
  margin: 0;
  /* 移除外边距 */
  border: none;
  /* 移除边框 */
  display: flex;
  justify-content: center;
  align-items: center;
}

.terminal {
  width: calc(100% - 20px);
  height: calc(100vh - 70px);
  padding: 0;
  /* 移除内边距 */
  margin: 0;
  /* 移除外边距 */
  border: none;
  /* 移除边框 */
}
</style>
<style>
/* 自定义 xterm 内部滚动条 */
.xterm .xterm-viewport {
  overflow-y: auto !important;
  /* 启用垂直滚动条 */
  overflow-x: hidden;
  /* 隐藏水平滚动条 */
}

/* 美化滚动条（可选） */
.xterm .xterm-viewport::-webkit-scrollbar {
  display: none;
}
</style>