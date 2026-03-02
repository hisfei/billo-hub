<template>
  <div class="manager-container">
    <div style="margin-bottom: 20px;">
      <el-button type="primary" icon="el-icon-plus" @click="showAdd = true">创建新 Agent</el-button>
      <el-button icon="el-icon-refresh" @click="fetchAgents">刷新状态</el-button>
    </div>

    <el-table :data="agentList" border style="width: 100%" v-loading="loading">
      <el-table-column prop="id" label="ID" width="150"></el-table-column>
      <el-table-column prop="name" label="名称"></el-table-column>
      <el-table-column label="运行状态" width="120">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status === 'RUNNING' ? 'success' : 'info'">
            {{ scope.row.status }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="lastActive" label="最后活跃时间" :formatter="formatTime"></el-table-column>
      <el-table-column label="操作" width="200">
        <template slot-scope="scope">
          <el-button size="mini" type="text" @click="handleRun(scope.row)">启动</el-button>
          <el-button size="mini" type="text" style="color: #F56C6C" @click="handleDelete(scope.row.id)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog title="创建 Agent" :visible.sync="showAdd" width="300px">
      <el-form :model="addForm">
        <el-form-item label="Agent 名称">
          <el-input v-model="addForm.name" placeholder="请输入名称"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="showAdd = false">取消</el-button>
        <el-button type="primary" @click="submitAdd">确定</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getAgents, createAgent, deleteAgent } from '@/api/agent';

export default {
  data() {
    return {
      agentList: [],
      loading: false,
      showAdd: false,
      addForm: { name: '' }
    }
  },
  mounted() {
    this.fetchAgents();
  },
  methods: {
    async fetchAgents() {
      this.loading = true;
      try {
        const res = await getAgents();
        this.agentList = res.data;
      } finally {
        this.loading = false;
      }
    },
    async submitAdd() {
      await createAgent(this.addForm);
      this.$message.success('创建成功');
      this.showAdd = false;
      this.fetchAgents();
    },
    async handleDelete(id) {
      await this.$confirm('确定要移除此 Agent 吗？');
      await deleteAgent(id);
      this.$message.success('已删除');
      this.fetchAgents();
    },
    handleRun(agent) {
      // 这里的逻辑可以跳转到之前的对话控制台，或者直接通过 WS 启动
      this.$emit('select-agent', agent.id);
    },
    formatTime(row) {
      return new Date(row.lastActive).toLocaleString();
    }
  }
}
</script>