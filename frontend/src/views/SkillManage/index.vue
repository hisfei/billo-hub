<template>
  <div class="skill-manage-container">
    <div class="manage-container">
      <div class="manage-header">
        <h3><i class="el-icon-s-tools"></i> {{ $t('skill.manage_title') }}</h3>
        <el-button type="primary" size="small" @click="showAddDialog = true" disabled>
          <i class="el-icon-plus"></i> {{ $t('skill.add_skill') }}
        </el-button>
      </div>
      <el-table :data="skillList" style="width: 100%"  :default-sort="{prop: 'id', order: 'descending'}">
        <el-table-column prop="id" :label="$t('skill.id')" width="180" sortable></el-table-column>
        <el-table-column prop="name" :label="$t('skill.name')" min-width="80"></el-table-column>
        <el-table-column prop="desc" :label="$t('skill.desc')" min-width="300"></el-table-column>
        <el-table-column :label="$t('skill.status')" width="120">
          <template slot-scope="scope">
            <el-tag :type="scope.row.status === 'ENABLED' ? 'success' : 'danger'">
              {{ scope.row.status === 'ENABLED' ? $t('skill.enabled') : $t('skill.disabled') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column :label="$t('agent.actions')" width="180">
          <template slot-scope="scope">
            <el-button type="primary" size="mini" @click="toggleSkillStatus(scope.row)">
              {{ scope.row.status === 'ENABLED' ? $t('skill.disable') : $t('skill.enable') }}
            </el-button>
            <el-button type="danger" size="mini" @click="deleteSkill(scope.row.id)" icon="el-icon-delete" style="margin-left: 8px;"></el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog :title="$t('skill.add_skill')" :visible.sync="showAddDialog" width="400px" :close-on-click-modal="false">
      <el-form :model="skillForm" :rules="skillRules" ref="skillForm">
        <el-form-item :label="$t('skill.name')" prop="name">
          <el-input v-model="skillForm.name" :placeholder="$t('skill.name_placeholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('skill.desc')" prop="desc">
          <el-input v-model="skillForm.desc" type="textarea" :rows="3" :placeholder="$t('skill.desc_placeholder')"></el-input>
        </el-form-item>
      </el-form>
      <span slot="footer">
        <el-button @click="showAddDialog = false; resetSkillForm()">{{ $t('llm.cancel') }}</el-button>
        <el-button type="primary" @click="addSkill">{{ $t('llm.confirm') }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { getSkillList, addSkill as apiAddSkill, toggleSkillStatus as apiToggleSkillStatus, deleteSkill as apiDeleteSkill } from '@/api/skill';

export default {
  name: "SkillManage",
  data() {
    return {
      skillList: [],
      showAddDialog: false,
      skillForm: {
        name: '',
        desc: ''
      },
      skillRules: {
        name: [{ required: true, message: this.$t('skill.name_required'), trigger: 'blur' }]
      }
    };
  },
  created() {
    this.loadSkillList();
  },
  methods: {
    async loadSkillList() {
      try {
        this.skillList = await getSkillList() || [];
      } catch (error) {
        this.$message.error(this.$t('skill.fetch_list_failed'));
        console.error('Failed to load skill list:', error);
      }
    },
    async addSkill() {
      this.$refs.skillForm.validate(async (valid) => {
        if (valid) {
          try {
            await apiAddSkill({ name: this.skillForm.name, desc: this.skillForm.desc || '无描述' });
            this.$message.success(this.$t('skill.add_success'));
            this.showAddDialog = false;
            this.resetSkillForm();
            this.loadSkillList();
          } catch (error) {
            this.$message.error(this.$t('skill.add_failed'));
            console.error('Failed to add skill:', error);
          }
        }
      });
    },
    async toggleSkillStatus(skill) {
      try {
        const newStatus = skill.status === 'ENABLED' ? 'DISABLED' : 'ENABLED';
        await apiToggleSkillStatus({ id: skill.id, status: newStatus });
        this.$message.success(this.$t('skill.status_update_success'));
        this.loadSkillList();
      } catch (error) {
        this.$message.error(this.$t('skill.status_update_failed'));
        console.error('Failed to toggle skill status:', error);
      }
    },
    async deleteSkill(id) {
      try {
        await this.$confirm(this.$t('skill.delete_confirm'), this.$t('llm.prompt'), { type: 'warning' });
        await apiDeleteSkill(id);
        this.$message.success(this.$t('skill.delete_success'));
        this.loadSkillList();
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error(this.$t('skill.delete_failed'));
        } else {
          this.$message.info(this.$t('skill.delete_cancel'));
        }
      }
    },
    resetSkillForm() {
      if (this.$refs.skillForm) {
        this.$refs.skillForm.resetFields();
      }
      this.skillForm = { name: '', desc: '' };
    }
  }
};
</script>

<style scoped>
.skill-manage-container {
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
  color: #E6A23C;
}
</style>