<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import ScanDialog from "./components/ScanDialog.vue";
import { EventUnpackProgress, ScanPathItem } from "./entries/entries";
import DataTable from 'primevue/datatable';
import Column from 'primevue/column';
import { formatSize, formatTime, UnpackStatusType, useAppToast } from "./entries/util";
import Toast from 'primevue/toast';
import UnpackDialog from "./components/UnpackDialog.vue";
import SettingsDialog from "./components/SettingsDialog.vue";
import {wechat} from "../wailsjs/go/models";
import WxapkgItem = wechat.WxapkgItem;
import UnpackOptions = wechat.UnpackOptions;
import * as AppService from '../wailsjs/go/main/AppService';

interface OutputDirSelectionPayload {
  outputDir: string;
  setAsDefault: boolean;
  useInSession: boolean;
}

interface StoredOutputDirectoryPreference {
  defaultOutputDir: string;
  lastOutputDir: string;
}

const OUTPUT_DIRECTORY_PREF_STORAGE_KEY = 'wxapkg:output-directory:preference:v1';

const scanDialogVisible = ref(false)
const unpackDialogVisible = ref(false)
const settingsDialogVisible = ref(false)
const search = ref<string>('')
const wxapkgItems = ref<WxapkgItem[]>([]);
const toast = useAppToast()
const version = ref<string>('v0.0.0')
const github = ref<string>('https://github.com')
const selectedWxapkgItem = ref<WxapkgItem | null>(null);
const tableKey = ref<string>('main-table')
const defaultOutputDir = ref('')
const lastOutputDir = ref('')
const sessionOutputDir = ref('')

const filteredItems = computed(() => {
  if (!search.value.trim()) return wxapkgItems.value
  const queryStr = search.value.toLowerCase().trim()
  return wxapkgItems.value.filter(item =>
    item.WxId.toLowerCase().includes(queryStr) ||
    item.Location.toLowerCase().includes(queryStr)
  )
})

function openUrl(url: string) {
  AppService.OpenUrl(url)
}

function openFolder(folder: string) {
  AppService.OpenPath(folder).catch(e => toast.error('打开目录失败', e))
}

function confirmScan(path: ScanPathItem) {
  AppService.ScanWxapkgItem(path.path, path.scan)
    .then((v: WxapkgItem[]) => {
      if (!v) v = []
      wxapkgItems.value = v
      toast.info('扫描小程序完成', `共 ${v.length} 个结果`)
    })
    .catch(e => toast.error('扫描小程序出错', e))
}

function copyPath(path: string) {
  AppService.ClipboardSetText(path)
    .then(() => toast.info('成功', '复制路径成功'))
    .catch(() => toast.error('失败', '复制路径失败'))
}

function unpack(item: WxapkgItem) {
  selectedWxapkgItem.value = item
  unpackDialogVisible.value = true
}

function confirmUnpack(options: UnpackOptions) {
  if (selectedWxapkgItem.value) {
    AppService.UnpackWxapkgItem(selectedWxapkgItem.value, options)
  }
}

function loadOutputDirectoryPreference() {
  try {
    const raw = localStorage.getItem(OUTPUT_DIRECTORY_PREF_STORAGE_KEY)
    if (!raw) return
    const parsed = JSON.parse(raw) as Partial<StoredOutputDirectoryPreference>
    defaultOutputDir.value = (parsed.defaultOutputDir ?? '').trim()
    lastOutputDir.value = (parsed.lastOutputDir ?? '').trim()
  } catch {
    defaultOutputDir.value = ''
    lastOutputDir.value = ''
  }
}

function persistOutputDirectoryPreference() {
  const data: StoredOutputDirectoryPreference = {
    defaultOutputDir: defaultOutputDir.value.trim(),
    lastOutputDir: lastOutputDir.value.trim(),
  }
  localStorage.setItem(OUTPUT_DIRECTORY_PREF_STORAGE_KEY, JSON.stringify(data))
}

function saveDefaultOutputDir(path: string) {
  defaultOutputDir.value = path.trim()
  if (sessionOutputDir.value && sessionOutputDir.value === defaultOutputDir.value) {
    sessionOutputDir.value = ''
  }
  persistOutputDirectoryPreference()
}

function handleOutputDirSelected(payload: OutputDirSelectionPayload) {
  const dir = payload.outputDir.trim()
  if (!dir) {
    return
  }

  lastOutputDir.value = dir

  if (payload.setAsDefault) {
    defaultOutputDir.value = dir
  }

  if (payload.useInSession) {
    sessionOutputDir.value = dir
  } else if (sessionOutputDir.value === dir) {
    sessionOutputDir.value = ''
  }

  persistOutputDirectoryPreference()
}

function openUnpackResultDirectory(path: string) {
  openFolder(path)
}

function handleDialogHide() {
  if (selectedWxapkgItem.value) {
    const index = wxapkgItems.value.findIndex(item => item.UUID === selectedWxapkgItem.value!.UUID)
    if (index !== -1) {
      wxapkgItems.value = wxapkgItems.value.map((item, i) =>
        i === index ? { ...item, ...selectedWxapkgItem.value! } : item
      )
    }
  }
  selectedWxapkgItem.value = null
}

function clearAll() {
  wxapkgItems.value = []
  toast.info('清空', '已清空所有小程序')
}

function getStatusIcon(status: string): string {
  const map: Record<string, string> = {
    [UnpackStatusType.Running]:  'pi pi-spin pi-spinner',
    [UnpackStatusType.Finished]: 'pi pi-folder-open',
    [UnpackStatusType.Error]:    'pi pi-times',
  }
  return map[status] ?? 'pi pi-box'
}

const notifiedUuids = new Set<string>()

function processProgress(uuid: string) {
  AppService.GetWxapkgItem(uuid).then(data => {
    if (!data) return
    const index = wxapkgItems.value.findIndex(item => item.UUID === uuid)
    if (index === -1) return

    const currentItem = wxapkgItems.value[index]
    const currentStatus = currentItem.UnpackStatus

    if ((currentStatus === UnpackStatusType.Finished || currentStatus === UnpackStatusType.Error)
        && data.UnpackStatus === UnpackStatusType.Running) {
      return
    }

    const updatedItem: WxapkgItem = { ...currentItem, ...data }
    wxapkgItems.value = wxapkgItems.value.map((item, i) => i === index ? updatedItem : item)

    if (unpackDialogVisible.value && selectedWxapkgItem.value?.UUID === uuid) {
      selectedWxapkgItem.value = updatedItem
    }

    if (!notifiedUuids.has(uuid)) {
      if (data.UnpackStatus === UnpackStatusType.Finished && !unpackDialogVisible.value) {
        notifiedUuids.add(uuid)
        toast.info('解包完成', `输出路径：${data.UnpackSavePath}`)
      } else if (data.UnpackStatus === UnpackStatusType.Error) {
        notifiedUuids.add(uuid)
        toast.error('解包失败', `${data.UnpackErrorMessage}`)
      }
    }
  })
}

onMounted(() => {
  loadOutputDirectoryPreference()
  window.runtime.EventsOn(EventUnpackProgress, (uuid: string) => {
    processProgress(uuid)
  })
  AppService.Version().then(v => version.value = v)
  AppService.Github().then(v => github.value = v)
})

onBeforeUnmount(() => {
  window.runtime.EventsOff(EventUnpackProgress)
})
</script>

<template>
  <div class="app-shell">
    <!-- Search -->
    <div class="search-bar">
      <div class="search-wrapper">
        <i class="pi pi-search search-icon"></i>
        <input
          v-model="search"
          class="search-input"
          type="text"
          placeholder="搜索小程序 ID 或路径"
        />
      </div>
      <div class="search-actions">
        <button class="btn-text" style="margin-right: 6px" @click="settingsDialogVisible = true">
          <i class="pi pi-cog" style="font-size:13px"></i>
          设置
        </button>
        <button
          v-if="wxapkgItems.length > 0"
          class="btn-text"
          style="margin-right: 20px"
          @click="clearAll"
        >
          <i class="pi pi-trash" style="font-size:13px"></i>
          清空列表
        </button>
        <button class="btn-primary" @click="scanDialogVisible = true">
          <i class="pi pi-folder-open"></i>
          扫描小程序
        </button>
      </div>
    </div>

    <!-- Table -->
    <div class="table-area">
      <div class="table-fill">
        <DataTable
          :value="filteredItems"
          data-key="UUID"
          sortField="LastModifyTime"
          :sortOrder="-1"
          scrollable
          scrollHeight="100%"
          tableStyle="width: 100%; table-layout: fixed;"
          :key="tableKey"
        >
          <colgroup>
            <col style="width: 170px">
            <col style="width: 180px">
            <col style="width: 100px">
            <col style="min-width: 200px">
            <col style="width: 72px">
          </colgroup>
        <template #empty>
          <div class="table-empty">
            {{ search ? `没有搜索到与 "${search}" 相关的小程序` : '没有小程序，请扫描或添加' }}
          </div>
        </template>

        <Column header="小程序 ID" field="WxId" style="width: 170px">
          <template #body="{ data }">
            <div class="mono" style="font-size:13px; white-space:nowrap">{{ data.WxId }}</div>
          </template>
        </Column>

        <Column header="小程序名称" field="WxId" style="width: 170px">
          <template #body="{ data }">
            <div class="mono" style="font-size:13px; white-space:nowrap">{{ data.Title }}</div>
          </template>
        </Column>

        <Column header="修改时间" field="LastModifyTime" :sortable="true" style="width: 180px">
          <template #body="{ data }">
            <div style="white-space:nowrap">{{ formatTime(data.LastModifyTime, false) }}</div>
          </template>
        </Column>

        <Column header="大小" field="Size" style="width: 100px" headerClass="col-right" bodyClass="col-right">
          <template #body="{ data }">
            <div class="mono" style="font-size:13px; white-space:nowrap">{{ formatSize(data.Size) }}</div>
          </template>
        </Column>

        <Column header="路径" field="Location" style="min-width: 200px">
          <template #body="{ data }">
            <div
              class="path-cell"
              v-tooltip.bottom="`点击复制\n${data.Location}`"
              @click="copyPath(data.Location)"
            >
              {{ data.Location }}
            </div>
          </template>
        </Column>

        <Column header="解包" style="width: 72px" headerClass="col-center" bodyClass="col-center">
          <template #body="{ data }">
            <button
              :class="[
                'status-dot',
                {
                  'running':    data.UnpackStatus === UnpackStatusType.Running,
                  'finished':  data.UnpackStatus === UnpackStatusType.Finished,
                  'error':     data.UnpackStatus === UnpackStatusType.Error,
                  'idle':      data.UnpackStatus === UnpackStatusType.Idle || !data.UnpackStatus,
                }
              ]"
              v-tooltip.top="
                data.UnpackStatus === UnpackStatusType.Running  ? `解包中 ${Math.round(data.UnpackProgress)}%` :
                data.UnpackStatus === UnpackStatusType.Finished ? '打开目录' :
                data.UnpackStatus === UnpackStatusType.Error    ? (data.UnpackErrorMessage || '解包失败') :
                '解包'"
              :disabled="data.UnpackStatus === UnpackStatusType.Running || data.UnpackStatus === UnpackStatusType.Error"
              @click="data.UnpackStatus === UnpackStatusType.Finished ? openFolder(data.UnpackSavePath) :
                      data.UnpackStatus === UnpackStatusType.Idle || !data.UnpackStatus ? unpack(data) : undefined"
            >
              <i :class="getStatusIcon(data.UnpackStatus)"></i>
            </button>
          </template>
        </Column>
      </DataTable>
      </div>
    </div>

    <!-- Footer -->
    <div class="footer">
      <span class="footer-disclaimer">仅供学习研究使用，请勿用于任何侵权或非法用途</span>
      <a class="footer-right" @click.prevent="openUrl(github)">
        <span class="footer-version">{{ version }}</span>
        <svg viewBox="0 0 24 24" fill="currentColor" xmlns="http://www.w3.org/2000/svg" width="13" height="13">
          <path d="M12 2C6.477 2 2 6.484 2 12.017c0 4.425 2.865 8.18 6.839 9.504.5.092.682-.217.682-.483 0-.237-.008-.868-.013-1.703-2.782.605-3.369-1.343-3.369-1.343-.454-1.158-1.11-1.466-1.11-1.466-.908-.62.069-.608.069-.608 1.003.07 1.531 1.032 1.531 1.032.892 1.53 2.341 1.088 2.91.832.092-.647.35-1.088.636-1.338-2.22-.253-4.555-1.113-4.555-4.951 0-1.093.39-1.988 1.029-2.688-.103-.253-.446-1.272.098-2.65 0 0 .84-.27 2.75 1.026A9.564 9.564 0 0112 6.844c.85.004 1.705.115 2.504.337 1.909-1.296 2.747-1.027 2.747-1.027.546 1.379.202 2.398.1 2.651.64.7 1.028 1.595 1.028 2.688 0 3.848-2.339 4.695-4.566 4.943.359.309.678.92.678 1.855 0 1.338-.012 2.419-.012 2.747 0 .268.18.58.688.482A10.019 10.019 0 0022 12.017C22 6.484 17.522 2 12 2z"/>
        </svg>
      </a>
    </div>

    <!-- Dialogs -->
    <ScanDialog v-model:visible="scanDialogVisible" @confirm="confirmScan"/>
    <UnpackDialog
      v-model:visible="unpackDialogVisible"
      v-model:item="selectedWxapkgItem"
      :default-output-dir="defaultOutputDir"
      :last-output-dir="lastOutputDir"
      :session-output-dir="sessionOutputDir"
      @after-hide="handleDialogHide"
      @confirm="confirmUnpack"
      @output-dir-selected="handleOutputDirSelected"
      @open-directory="openUnpackResultDirectory"
    />
    <SettingsDialog
      v-model:visible="settingsDialogVisible"
      :default-output-dir="defaultOutputDir"
      @save="saveDefaultOutputDir"
    />
    <Toast position="bottom-right" />
  </div>
</template>

<style scoped>
.app-shell {
  display: flex;
  flex-direction: column;
  height: 100vh;
  background: var(--color-white);
  overflow: hidden;
}

/* Search */
.search-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 14px 20px 12px;
  flex-shrink: 0;
  border-bottom: 1px solid rgba(0, 0, 0, 0.06);
}

.search-actions {
  display: flex;
  align-items: center;
  gap: 4px;
  flex-shrink: 0;
}

/* Table */
.table-area {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
  padding: 14px 20px;
  min-height: 0;
}

.table-fill {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0;
  overflow: hidden;
  height: 100%;
}

.table-fill :deep(.p-datatable) {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.table-fill :deep(.p-datatable-table) {
  flex: 1;
}

.table-fill :deep(.p-datatable-thead) {
  flex-shrink: 0;
}

.table-fill :deep(.p-datatable-wrapper) {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.table-fill :deep(.p-datatable-tbody) {
  min-height: 100%;
}

.table-area :deep(.p-datatable-thead > tr > th) {
  background: #fafafa !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.08) !important;
  color: var(--color-text-tertiary) !important;
  font-size: 11px !important;
  font-weight: 600 !important;
  text-transform: uppercase !important;
  letter-spacing: 0.05em !important;
  padding: 10px 16px !important;
}

.table-area :deep(.p-datatable-thead > tr > th:last-child) {
  padding-right: 16px !important;
}

.table-area :deep(.p-sortable-column:not(:last-child)) {
  padding-right: 24px !important;
}

.table-area :deep(.p-datatable-tbody > tr > td) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
  padding: 11px 16px !important;
  color: var(--color-text-secondary) !important;
  font-size: 14px !important;
}

.table-area :deep(.p-datatable-tbody > tr:last-child > td) {
  border-bottom: none !important;
}

.table-area :deep(.p-datatable-tbody > tr:hover > td) {
  background: rgba(0, 0, 0, 0.018) !important;
}

.table-area :deep(.p-datatable-column-sorted) {
  background: transparent !important;
}

.table-empty {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 160px;
  color: var(--color-text-tertiary);
  font-size: 14px;
}

.status-dot i {
  font-size: 15px;
}

.path-cell {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  direction: rtl;
  text-align: left;
  font-family: "JetBrains Mono", "Cascadia Code", "Consolas", monospace;
  font-size: 13px;
  cursor: default;
}

/* Footer */
.footer {
  display: flex;
  align-items: center;
  justify-content: flex-end;
  padding: 8px 20px;
  border-top: 1px solid rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
  position: relative;
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 6px;
  /* background: var(--color-light-gray); */
  border-radius: 980px;
  padding: 5px 10px;
  flex-shrink: 0;
  cursor: pointer;
  text-decoration: none;
  /* color: rgba(0, 0, 0, 0.5); */
  transition: background 0.15s, color 0.15s;
}
.footer-right:hover {
  background: rgba(0, 0, 0, 0.1);
  /* color: var(--color-apple-blue); */
}

.footer-disclaimer {
  position: absolute;
  left: 50%;
  transform: translateX(-50%);
  font-size: 11px;
  color: rgba(0, 0, 0, 0.5);
  white-space: nowrap;
}

.footer-version {
  font-size: 11px;
  color: rgba(0, 0, 0, 0.5);
}

</style>
