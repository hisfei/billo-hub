<template>
  <div class="agent-manage-container">
    <div class="manage-container">
      <div class="manage-header">
        <h3><i class="el-icon-user-solid"></i> {{ $t('agent.manage_title') }}</h3>
        <el-button type="primary" size="small" @click="openAddDialog">
          <i class="el-icon-plus"></i> {{ $t('agent.add_agent') }}
        </el-button>
      </div>
      <el-table :data="agentList" style="width: 100%">
        <el-table-column prop="name" :label="$t('agent.name')" />
        <el-table-column prop="llm" :label="$t('agent.model')" />
        <el-table-column prop="persona" :label="$t('agent.persona')" min-width="160" show-overflow-tooltip />
        <el-table-column :label="$t('agent.skills')" min-width="200">
          <template slot-scope="scope">
            <span v-for="(s, i) in scope.row.skills" :key="i" class="skill-tag">{{ s }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="maxLoops" :label="$t('agent.max_loops')" />
        <el-table-column :label="$t('agent.surfing_status')">
          <template slot-scope="scope">
            <el-switch v-model="scope.row.openBackgroundSurfing" @change="handleStatusChange(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('agent.active_status')">
          <template slot-scope="scope">
            <el-switch v-model="scope.row.isActive" @change="handleStatusChange(scope.row)" />
          </template>
        </el-table-column>
        <el-table-column :label="$t('agent.actions')" width="100">
          <template slot-scope="scope">
            <el-button type="danger" size="mini" @click="deleteAgent(scope.row.id)" icon="el-icon-delete" />
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog :title="$t('agent.add_agent')" :visible.sync="showAddDialog" width="520px" :close-on-click-modal="false">
      <el-form :model="agentForm" :rules="agentRules" ref="agentForm" label-width="100px">
        <el-form-item :label="$t('agent.name')" prop="name">
          <el-input v-model="agentForm.name" :placeholder="$t('agent.name_placeholder')" />
        </el-form-item>
        <el-form-item :label="$t('agent.persona')" prop="persona">
          <el-input v-model="agentForm.persona" type="textarea" :rows="3" :placeholder="$t('agent.persona_placeholder')" />
        </el-form-item>
        <el-form-item :label="$t('agent.model')" prop="llm">
          <el-select v-model="agentForm.llm" :placeholder="$t('agent.model_placeholder')" style="width: 100%">
            <el-option v-for="i in llmOptions" :key="i.name" :label="i.name" :value="i.name" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('agent.active_status')" prop="isActive">
          <el-select v-model="agentForm.isActive" :placeholder="$t('agent.status_placeholder')" style="width: 100%">
            <el-option v-for="bo in boolOptions" :key="bo.id" :label="bo.name" :value="bo.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('agent.surfing_status')" prop="openBackgroundSurfing">
          <el-select v-model="agentForm.openBackgroundSurfing" :placeholder="$t('agent.status_placeholder')" style="width: 100%">
            <el-option v-for="bo in boolOptions" :key="bo.id" :label="bo.name" :value="bo.id" />
          </el-select>
        </el-form-item>
        <el-form-item :label="$t('agent.skills')" prop="skills">
          <el-select v-model="agentForm.skills" multiple :placeholder="$t('agent.skills_placeholder')" style="width: 100%">
            <el-option v-for="skill in skillOptions" :key="skill.id" :label="skill.name" :value="skill.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <div slot="footer" class="dialog-footer">
        <el-button @click="closeAddDialog">{{ $t('llm.cancel') }}</el-button>
        <el-button type="primary" @click="handleAddAgent">{{ $t('llm.confirm') }}</el-button>
      </div>
    </el-dialog>
  </div>
</template>

<script>
import { getAgents, createAgent, deleteAgent, updateAgentStatus, getLLMList } from "@/api/agent";
import { getSkillList } from "@/api/skill";

export default {
  name: "AgentManage",
  data() {
    return {
      agentList: [],
      showAddDialog: false,
      agentForm: {
        name: "",
        persona: "",
        skills: [],
        llm: [],
        isActive: true,
        maxLoops: 500,
        openBackgroundSurfing: false,
      },
      agentRules: {
        name: [{ required: true, message: this.$t('agent.name_required'), trigger: "blur" }],
        persona: [{ required: true, message: this.$t('agent.persona_required'), trigger: "blur" }],
        llm: [{ required: true, message: this.$t('agent.model_required'), trigger: "change" }],
        skills: [{ required: true, message: this.$t('agent.skills_required'), trigger: "change" }]
      },
      llmOptions: [],
      boolOptions: [
        { id: true, name: this.$t('agent.active') },
        { id: false, name: this.$t('agent.inactive') }
      ],
      skillOptions: []
    };
  },
  created() {
    this.getSkillList();
    this.getLLMList();
    this.fetchAgentList();
  },
  methods: {
    async getLLMList() {
      try {
        this.llmOptions = await getLLMList() || [];
      } catch (error) {
        console.error('Failed to load LLM list:', error);
        this.llmOptions = [];
      }
    },
    async getSkillList() {
      try {
        this.skillOptions = await getSkillList() || [];
      } catch (error) {
        console.error('Failed to load skill list:', error);
        this.skillOptions = [];
      }
    },
    openAddDialog() {
      this.showAddDialog = true;
    },
    closeAddDialog() {
      this.showAddDialog = false;
      this.$nextTick(() => {
        if (this.$refs.agentForm) this.$refs.agentForm.resetFields();
      });
    },
    async fetchAgentList() {
      try {
        this.agentList = await getAgents() || [];
      } catch (err) {
        this.$message.error(this.$t('agent.fetch_list_failed'));
      }
    },
    async handleAddAgent() {
      this.$refs.agentForm.validate(async valid => {
        if (!valid) return;
        try {
          await createAgent(this.agentForm);
          this.$message.success(this.$t('agent.add_success'));
          this.closeAddDialog();
          this.fetchAgentList();
        } catch (err) {
          this.$message.error(this.$t('agent.add_failed'));
        }
      });
    },
    async handleStatusChange(agent) {
      try {
        await updateAgentStatus({ id: agent.id, is_active: agent.is_active });
        this.$message.success(this.$t('agent.status_update_success'));
      } catch (err) {
        agent.is_active = !agent.is_active;
        this.$message.error(this.$t('agent.status_update_failed'));
      }
    },
    async deleteAgent(id) {
      try {
        await this.$confirm(this.$t('agent.delete_confirm'), this.$t('llm.prompt'), { type: "warning" });
        await deleteAgent({ id: id });
        this.$message.success(this.$t('agent.delete_success'));
        this.fetchAgentList();
      } catch (err) {
        if (err !== 'cancel') {
          this.$message.error(this.$t('agent.delete_failed'));
        } else {
          this.$message.info(this.$t('agent.delete_cancel'));
        }
      }
    }
  }
};
</script>

<style scoped>
.agent-manage-container {
  padding: 20px;
}
.manage-container {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}
.manage-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}
.manage-header h3 {
  margin: 0;
  font-size: 16px;
  display: flex;
  align-items: center;
}
.manage-header h3 i {
  margin-right: 8px;
  color: #67C23A;
}
.skill-tag {
  display: inline-block;
  background: #f4f4f5;
  border-radius: 4px;
  padding: 2px 6px;
  margin: 0 4px 4px 0;
  font-size: 12px;
}
.dialog-footer {
  text-align: right;
}
</style>