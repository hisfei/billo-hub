<template>
  <div class="llm-manage-container">
    <div class="manage-container">
      <div class="manage-header">
        <h3><i class="el-icon-s-operation"></i> {{ $t('llm.manage_title') }}</h3>
        <el-button type="primary" size="small" @click="dialogVisible = true">
          <i class="el-icon-plus"></i> {{ $t('llm.add_model') }}
        </el-button>
      </div>
      <el-table :data="llmList" style="width: 100%">
        <el-table-column prop="name" :label="$t('llm.model_name')"></el-table-column>
        <el-table-column prop="url" :label="$t('llm.request_url')"></el-table-column>
        <el-table-column prop="supportContextId" :label="$t('llm.support_context_id')">
          <template slot-scope="scope">
            <el-tag :type="scope.row.supportContextId ? 'success' : 'info'">
              {{ scope.row.supportContextId ? $t('llm.yes') : $t('llm.no') }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="contextExpire" :label="$t('llm.context_expire')"></el-table-column>
        <el-table-column :label="$t('agent.actions')" width="100">
          <template slot-scope="scope">
            <el-button type="danger" size="mini" @click="deleteLlm(scope.row.id)" icon="el-icon-delete"></el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog :title="$t('llm.add_model')" :visible.sync="dialogVisible" width="50%">
      <el-form :model="llmForm" :rules="llmRules" ref="llmForm" label-width="150px">
        <el-form-item :label="$t('llm.model_name')" prop="name">
          <el-input v-model="llmForm.name" :placeholder="$t('llm.model_name_placeholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('llm.request_url')" prop="url">
          <el-input v-model="llmForm.url" :placeholder="$t('llm.request_url_placeholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('llm.api_key')" prop="apiKey">
          <el-input v-model="llmForm.apiKey" :placeholder="$t('llm.api_key_placeholder')"></el-input>
        </el-form-item>
        <el-form-item :label="$t('llm.support_context_id')" prop="supportContextId">
          <el-switch v-model="llmForm.supportContextId"></el-switch>
        </el-form-item>
        <el-form-item :label="$t('llm.context_expire')" prop="contextExpire">
          <el-input-number v-model="llmForm.contextExpire" :min="0"></el-input-number>
          <span class="tip">{{ $t('llm.context_expire_tip') }}</span>
        </el-form-item>
      </el-form>
      <span slot="footer" class="dialog-footer">
        <el-button @click="dialogVisible = false">{{ $t('llm.cancel') }}</el-button>
        <el-button type="primary" @click="submitForm('llmForm')">{{ $t('llm.confirm') }}</el-button>
      </span>
    </el-dialog>
  </div>
</template>

<script>
import { addLlm, getLlmList, deleteLlm } from '@/api/llm';

export default {
  name: 'LlmManage',
  data() {
    return {
      llmList: [],
      dialogVisible: false,
      llmForm: {
        name: '',
        url: '',
        apiKey: '',
        supportContextId: false,
        contextExpire: 0
      },
      llmRules: {
        name: [
          { required: true, message: this.$t('llm.model_name_required'), trigger: 'blur' }
        ],
        url: [
          { required: true, message: this.$t('llm.request_url_required'), trigger: 'blur' },
          { type: 'url', message: this.$t('llm.request_url_invalid'), trigger: ['blur', 'change'] }
        ],
        apiKey: [
          { required: true, message: this.$t('llm.api_key_required'), trigger: 'blur' }
        ]
      }
    };
  },
  created() {
    this.fetchLlmList();
  },
  methods: {
    async fetchLlmList() {
      try {
        this.llmList = await getLlmList();
      } catch (error) {
        this.$message.error(this.$t('llm.fetch_list_failed'));
        console.error('Fetch LLM list failed:', error);
      }
    },
    submitForm(formName) {
      this.$refs[formName].validate(async (valid) => {
        if (valid) {
          try {
            await addLlm(this.llmForm);
            this.$message.success(this.$t('llm.add_model_success'));
            this.dialogVisible = false;
            this.resetForm(formName);
            this.fetchLlmList();
          } catch (error) {
            this.$message.error(this.$t('llm.add_model_failed'));
            console.error('Add LLM failed:', error);
          }
        }
      });
    },
    resetForm(formName) {
      if (this.$refs[formName]) {
        this.$refs[formName].resetFields();
      }
    },
    async deleteLlm(id) {
      try {
        await this.$confirm(this.$t('llm.delete_confirm'), this.$t('llm.prompt'), { type: 'warning' });
        await deleteLlm(id);
        this.$message.success(this.$t('llm.delete_success'));
        this.fetchLlmList();
      } catch (error) {
        if (error !== 'cancel') {
          this.$message.error(this.$t('llm.delete_failed'));
        } else {
          this.$message.info(this.$t('llm.delete_cancel'));
        }
      }
    }
  }
};
</script>

<style scoped>
.llm-manage-container {
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
  color: #409EFF;
}
.tip {
  margin-left: 10px;
  font-size: 12px;
  color: #909399;
}
.dialog-footer {
  text-align: right;
}
</style>