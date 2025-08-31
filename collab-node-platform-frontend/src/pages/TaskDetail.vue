<template>
  <div class="task-detail-page">
    <el-row justify="space-between" align="middle" style="margin-bottom: 20px;">
      <el-col>
        <h2 style="display:inline-block;">{{ taskName }}</h2>
      </el-col>
      <el-col>
        <el-button :disabled="!selectedNode || taskStatus===2" @click="showCreate = true" type="primary">新建子节点</el-button>
        <el-button
          v-if="selectedNode"
          type="danger"
          style="margin-left:10px;"
          :disabled="taskStatus===2"
          @click="onDeleteNode"
        >删除节点</el-button>
        <el-input v-model="searchKey" placeholder="分支关键词" style="width:200px;margin-left:10px;" @keyup.enter="onSearchAndLocate" />
        <el-button style="margin-left:5px;" @click="onSearchAndLocate" type="primary">定位</el-button>
        <el-button style="margin-left:5px;" @click="goToRootNode" type="info">回到根节点</el-button>
        <el-button style="margin-left:5px;" :disabled="!selectedNode" @click="onFocusSelected" type="success">专注当前节点分支(节点点击)</el-button>
        <el-button style="margin-left:10px;" @click="logout">退出登录</el-button>
        <el-button type="info" @click="goBackList" style="margin-left:10px;font-weight:bold;font-size:16px;border-radius:24px;padding:8px 28px;background:#34495e;color:#fff;">返回任务列表</el-button>
      </el-col>
    </el-row>
      <el-row style="margin-bottom:10px;">
        <el-col>
          <el-button v-if="focusedNodeId" @click="clearFocus" size="small">恢复全局视图</el-button>
        </el-col>
      </el-row>
  <div id="cy" style="height:600px;border:1px solid #eee;"></div>
    <el-dialog v-model="showCreate" title="新建子节点" @open="onOpenCreateDialog">
      <el-form :model="createForm" :disabled="taskStatus===2">
        <el-form-item label="父节点">
          <el-input v-model="parentNodeName" disabled />
        </el-form-item>
        <el-form-item label="节点名称">
          <el-input v-model="createForm.nodeName" />
        </el-form-item>
        <el-form-item label="节点内容">
          <el-input v-model="createForm.nodeContent" type="textarea" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreate = false">取消</el-button>
        <el-button type="primary" @click="onCreateNode">创建</el-button>
      </template>
    </el-dialog>
    <el-dialog v-model="showDetailDialog" title="节点详情" width="500px">
      <el-form :model="detailForm" :disabled="taskStatus===2">
        <el-form-item label="节点名称">
          <el-input v-model="detailForm.nodeName" />
        </el-form-item>
        <el-form-item label="节点内容">
          <el-input v-model="detailForm.nodeContent" type="textarea" />
        </el-form-item>
        <el-form-item label="站点">
          <el-input v-model="detailForm.site" />
        </el-form-item>
        <el-form-item label="结果">
          <el-input v-model="detailForm.result" />
        </el-form-item>
        <el-form-item label="下一步">
          <el-input v-model="detailForm.nextStep" />
        </el-form-item>
        <el-form-item label="版本">
          <el-input v-model="detailForm.version" disabled />
        </el-form-item>
        <el-form-item label="创建者">
          <el-input v-model="detailForm.creatorName" disabled />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDetailDialog = false">取消</el-button>
        <el-button type="primary" @click="onSaveDetail" :disabled="taskStatus===2">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>
<script setup>
// 搜索并自动定位到第一个匹配节点
function onSearchAndLocate() {
  if (!searchKey.value || !cy) return
  // 搜索所有节点名/内容
  const found = cy.nodes().filter(n => {
    const raw = n.data('raw') || {}
    return (raw.nodeName && raw.nodeName.includes(searchKey.value)) || (raw.nodeContent && raw.nodeContent.includes(searchKey.value))
  })
  if (found.length > 0) {
    const node = found[0]
    focusBranch(node.id())
    cy.center(node)
    cy.$(':selected').unselect()
    node.select()
    selectedNode.value = node.data('raw')
    ElMessage.success('已定位到第一个匹配节点')
  } else {
    ElMessage.warning('未找到匹配节点')
  }
}

// 一键回到根节点
function goToRootNode() {
  if (!cy) return
  // 假设根节点parentNodeId为null
  const root = cy.nodes().filter(n => !n.data('raw').parentNodeId)[0]
  if (root) {
    focusBranch(root.id())
    cy.center(root)
    cy.$(':selected').unselect()
    root.select()
    selectedNode.value = root.data('raw')
    ElMessage.success('已回到根节点')
  } else {
    ElMessage.warning('未找到根节点')
  }
}
import { nextTick } from 'vue'
// 专注分支相关
const focusedNodeId = ref(null)
function focusBranch(nodeId) {
  focusedNodeId.value = nodeId
  if (!cy) return
  // 1. 找到从根到该节点的路径
  const path = []
  let current = cy.getElementById(nodeId)
  while (current && current.length && current.data('raw') && current.data('raw').parentNodeId) {
    path.unshift(current.id())
    current = cy.getElementById(current.data('raw').parentNodeId)
  }
  if (current && current.length) path.unshift(current.id())
  // 2. 高亮该分支，其他变灰
  cy.nodes().forEach(n => {
    if (path.includes(n.id())) {
      n.addClass('focus-branch')
      n.removeClass('blur-branch')
    } else {
      n.removeClass('focus-branch')
      n.addClass('blur-branch')
    }
  })
  cy.edges().forEach(e => {
    if (path.includes(e.source().id()) && path.includes(e.target().id())) {
      e.addClass('focus-branch')
      e.removeClass('blur-branch')
    } else {
      e.removeClass('focus-branch')
      e.addClass('blur-branch')
    }
  })
}
function clearFocus() {
  focusedNodeId.value = null
  if (!cy) return
  cy.nodes().forEach(node => {
    node.style('display', 'element')
    node.removeClass('focus-branch blur-branch')
  })
  cy.edges().forEach(edge => {
    edge.style('display', 'element')
    edge.removeClass('focus-branch blur-branch')
  })
}
const onDeleteNode = async () => {
  if (taskStatus.value === 2) {
    ElMessage.warning('已结束任务不允许删除节点')
    return
  }
  if (!selectedNode.value) return
  if (!selectedNode.value.parentNodeId) {
    ElMessage.warning('根节点不允许删除')
    return
  }
  // 直接请求后端批量删除（后端已递归删除所有子节点）
  if (!confirm('确定要删除该节点及其所有子节点吗？')) return
  try {
    await deleteNode(selectedNode.value.nodeId)
    ElMessage.success('删除成功')
    selectedNode.value = null
    fetchNodes()
  } catch (e) {}
}
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getNodeList, createNode, editNode, searchNode, deleteNode, getTaskList } from '../api'
import { connectWS, disconnectWS } from '../utils/ws'
import { ElMessage } from 'element-plus'
import cytoscape from 'cytoscape'
import { useUserStore } from '../store/user'


const showDetailDialog = ref(false)
const detailForm = ref({
  nodeName: '',
  nodeContent: '',
  site: '',
  result: '',
  nextStep: '',
  version: '',
  creatorName: ''
})

function showDetail(node) {
  detailForm.value = {
    nodeName: node.nodeName || '',
    nodeContent: node.nodeContent || '',
    site: node.site || '',
    result: node.result || '',
    nextStep: node.nextStep || '',
    version: node.version || '',
    creatorName: node.creatorUsername || node.creatorId || ''
  }
  showDetailDialog.value = true
}

const userStore = useUserStore()
const router = useRouter()
const logout = () => {
  userStore.logout()
  router.replace('/login')
}

const route = useRoute()
const taskId = route.params.id
const taskName = ref('')
const showCreate = ref(false)
const createForm = ref({ nodeName: '', nodeContent: '' })
const parentNodeName = ref('')
const searchKey = ref('')
let cy = null
const selectedNode = ref(null)
const taskStatus = ref(1) // 默认1:开启，2:已结束

const fetchNodes = async () => {
  const res = await getNodeList(taskId)
  renderGraph(res)
  // 获取任务状态
  try {
    const tasks = await getTaskList()
    const t = tasks.find(t => t.taskId == taskId)
    if (t) taskStatus.value = t.status
  } catch {}
}

const renderGraph = (nodes) => {
  // 记录所有节点id，便于判断末尾节点
  const allNodeIds = new Set(nodes.map(n => n.nodeId))
  const parentIds = new Set(nodes.filter(n => n.parentNodeId).map(n => n.parentNodeId))
  // 末尾节点id集合
  const leafNodeIds = new Set([...allNodeIds].filter(id => !parentIds.has(id)))
  window._leafNodeIds = leafNodeIds
  if (!cy) {
    cy = cytoscape({
      container: document.getElementById('cy'),
      style: [
        {
          selector: 'node',
          style: {
            'background-color': '#e6f7ff',
            'label': 'data(label)',
            'text-valign': 'center',
            'text-halign': 'center',
            'width': 140,
            'height': 48,
            'font-size': 13,
            'border-width': 1,
            'border-color': '#333',
            'color': '#333',
            'text-wrap': 'wrap',
            'text-max-width': 130
          }
        },
        {
          selector: 'edge',
          style: {
            'width': 2,
            'line-color': '#888',
            'target-arrow-color': '#888',
            'target-arrow-shape': 'triangle'
          }
        },
        {
          selector: 'node:selected',
          style: {
            'background-color': '#ffd666',
            'border-color': '#ff9900'
          }
        }
      ],
      userZoomingEnabled: true,
      userPanningEnabled: true,
      boxSelectionEnabled: false,
      autoungrabify: true // 禁止节点拖动
    })
    cy.on('tap', 'node', evt => {
      const node = evt.target.data('raw')
      selectedNode.value = node
    })
    cy.on('dbltap', 'node', evt => {
      const node = evt.target.data('raw')
      showDetail(node)
    })
    cy.on('tap', evt => {
      if (evt.target === cy) selectedNode.value = null
    })
    // 右键菜单：专注分支
    cy.on('cxttap', 'node', evt => {
      evt.preventDefault()
      const nodeId = evt.target.id()
      focusBranch(nodeId)
    })
  } else {
    cy.elements().remove()
  }
  // 构建cytoscape数据
  const cyNodes = nodes.map(n => ({
    data: {
      id: n.nodeId,
      label: `${n.nodeName}\n创建者:${n.creatorUsername || n.creatorId}`,
      raw: n
    }
  }))
  const cyEdges = nodes.filter(n => n.parentNodeId).map(n => ({
    data: { source: n.parentNodeId, target: n.nodeId }
  }))
  cy.add([...cyNodes, ...cyEdges])
  cy.layout({ name: 'breadthfirst', directed: true, padding: 40, spacingFactor: 1.2 }).run()
  // 渲染后自动专注/全局
  nextTick(() => {
    if (focusedNodeId.value) focusBranch(focusedNodeId.value)
    else clearFocus()
  })
}
    const onCreateNode = async () => {
      if (taskStatus.value === 2) {
        ElMessage.warning('已结束任务不允许新增节点')
        return
      }
      if (!createForm.value.nodeName) {
        ElMessage.error('请输入节点名称')
        return
      }
      if (!selectedNode.value) {
        ElMessage.error('请先选中父节点')
        showCreate.value = false
        return
      }
      // 检查全局是否有同名节点
      let hasSameName = false
      if (cy) {
        cy.nodes().forEach(node => {
          const raw = node.data('raw')
          if (raw && raw.nodeName === createForm.value.nodeName) {
            hasSameName = true
          }
        })
      }
      if (hasSameName) {
        ElMessage.error('全局已存在同名节点')
        return
      }
      // 新节点位置：父节点下方偏移
      let pos = { x: 100, y: 100 }
      if (selectedNode.value.site) {
        try { pos = JSON.parse(selectedNode.value.site) } catch {}
      }
      const newPos = { x: pos.x, y: pos.y + 100 }
      try {
        await createNode({
          taskId,
          parentNodeId: selectedNode.value.nodeId,
          nodeName: createForm.value.nodeName,
          nodeContent: createForm.value.nodeContent,
          site: JSON.stringify(newPos)
        })
        ElMessage.success('创建成功')
        showCreate.value = false
        createForm.value = { nodeName: '', nodeContent: '' }
        fetchNodes()
      } catch (e) {}
    }

    const onSearch = async () => {
      if (!searchKey.value) return
      const res = await searchNode(taskId, searchKey.value)
      if (res.length) {
        // 专注第一个匹配节点及其所有子节点
        const n = res[0]
        if (cy) {
          // 递归获取所有子节点id
          const getAllDescendants = (id) => {
            const descendants = [id]
            cy.nodes().forEach(node => {
              if (node.data('raw') && node.data('raw').parentNodeId === id) {
                descendants.push(...getAllDescendants(node.id()))
              }
            })
            return descendants
          }
          const focusIds = getAllDescendants(n.nodeId)
          cy.nodes().forEach(node => {
            if (focusIds.includes(node.id())) {
              node.style('display', 'element')
              node.addClass('focus-branch')
              node.removeClass('blur-branch')
            } else {
              node.style('display', 'none')
              node.removeClass('focus-branch')
              node.removeClass('blur-branch')
            }
          })
          cy.edges().forEach(edge => {
            if (focusIds.includes(edge.source().id()) && focusIds.includes(edge.target().id())) {
              edge.style('display', 'element')
              edge.addClass('focus-branch')
              edge.removeClass('blur-branch')
            } else {
              edge.style('display', 'none')
              edge.removeClass('focus-branch')
              edge.removeClass('blur-branch')
            }
          })
          focusedNodeId.value = n.nodeId
        }
      } else {
        ElMessage.warning('未找到相关节点')
      }
    }

    // 专注当前选中节点分支
    function onFocusSelected() {
      if (selectedNode.value && selectedNode.value.nodeId) {
        // 递归获取所有子节点id
        const getAllDescendants = (id) => {
          const descendants = [id]
          cy.nodes().forEach(node => {
            if (node.data('raw') && node.data('raw').parentNodeId === id) {
              descendants.push(...getAllDescendants(node.id()))
            }
          })
          return descendants
        }
        const focusIds = getAllDescendants(selectedNode.value.nodeId)
        cy.nodes().forEach(node => {
          if (focusIds.includes(node.id())) {
            node.style('display', 'element')
            node.addClass('focus-branch')
            node.removeClass('blur-branch')
          } else {
            node.style('display', 'none')
            node.removeClass('focus-branch')
            node.removeClass('blur-branch')
          }
        })
        cy.edges().forEach(edge => {
          if (focusIds.includes(edge.source().id()) && focusIds.includes(edge.target().id())) {
            edge.style('display', 'element')
            edge.addClass('focus-branch')
            edge.removeClass('blur-branch')
          } else {
            edge.style('display', 'none')
            edge.removeClass('focus-branch')
            edge.removeClass('blur-branch')
          }
        })
        focusedNodeId.value = selectedNode.value.nodeId
      }
    }

    onMounted(() => {
      fetchNodes()
      connectWS(taskId, (event, data) => {
        if (["newNode", "editNode", "deleteNode", "updateTask", "deleteTask"].includes(event)) {
          fetchNodes()
        }
        if (event === "updateTask" && data && data.status !== undefined) {
          // 可选：根据任务状态做额外处理
        }
        if (event === "deleteTask") {
          // 任务被删除，自动返回任务列表页
          ElMessage.warning('该任务已被删除')
          router.replace('/')
        }
      })
    })
    onUnmounted(() => {
      disconnectWS()
    })
    const goBackList = () => {
      router.replace('/')
    }
    </script>
<style scoped>
.task-detail-page {
  position: fixed;
  inset: 0;
  width: 100vw;
  height: 100vh;
  background: linear-gradient(135deg, #e3f2fd 0%, #fff 100%);
  display: flex;
  flex-direction: column;
  padding: 0;
  margin: 0;
  border-radius: 0;
  box-shadow: none;
  z-index: 1;
}
#cy {
  flex: 1;
  min-height: 0;
  width: 100%;
  height: 100%;
  background: linear-gradient(120deg, #f5faff 0%, #e1f5fe 100%);
  border-radius: 12px;
  box-shadow: 0 4px 24px #b3e5fc44;
}
/**** cytoscape节点美化（通过class） ****/
.cytoscape-node {
  border-radius: 18px;
  background: linear-gradient(135deg, #90caf9 0%, #e3f2fd 100%);
  box-shadow: 0 2px 12px #90caf988;
  transition: box-shadow 0.2s, background 0.2s;
}
.cytoscape-node:hover {
  background: linear-gradient(135deg, #42a5f5 0%, #b3e5fc 100%);
  box-shadow: 0 4px 24px #42a5f5cc;
}
.cytoscape-edge {
  stroke: url(#edge-gradient);
  stroke-width: 3px;
}
/**** 分支专注更明显 ****/
.focus-branch {
  background: #fffde7 !important;
  border-color: #ff9800 !important;
  color: #e65100 !important;
  box-shadow: 0 0 16px #ff980088 !important;
  z-index: 10;
}
.blur-branch {
  opacity: 0.12 !important;
  filter: grayscale(0.7);
}
/**** 选中节点高亮 ****/
.cytoscape-node.selected, .focus-branch.selected {
  background: #ffe082 !important;
  border-color: #ffb300 !important;
  box-shadow: 0 0 24px #ffd54f !important;
}
</style>
