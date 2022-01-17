<template>
  <el-dialog
    :title="title"
    :visible.sync="visible"
    width="500px"
    :before-close="handleClose"
  >
    <el-form
      ref="formData"
      :rules="rules"
      :model="formData"
      label-width="100px"
      label-position="right"
      style="width: 400px"
      status-icon
    >
      <el-form-item label="名称：" prop="name">
        <el-input v-model="formData.name" />
      </el-form-item>
      <el-form-item label="指纹内容" prop="fingerprint">
        <el-input v-model="formData.fingerprint" />
      </el-form-item>
      <el-form-item label="备注：" prop="remark">
        <el-input v-model="formData.remark" type="textarea" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" size="mini" @click="submitForm('formData')">确定</el-button>
        <el-button size="mini" @click="handleClose">取消</el-button>
      </el-form-item>
    </el-form>
  </el-dialog>
</template>
<script>
import api from '@/api/fingerprint'
export default {
  props: {
    title: {
      type: String,
      default: ''
    },
    visible: {
      type: Boolean,
      default: false
    },
    formData: {
      type: Object,
      default: () => ({})
    },
    remoteClose: {
      type: Function,
      default: () => () => {}
    }
  },

  data() {
    return {
      rules: {
        name: [
          { required: true, message: '请输入规则名称', trigger: 'blur' }
        ],
        fingerprint: [
          { required: true, message: '请选择规则类型', trigger: 'change' }
        ]
      }
    }
  },

  methods: {
    handleClose() {
      this.$refs['formData'].resetFields()
      this.remoteClose()
    },

    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          this.submitData()
        } else {
          return false
        }
      })
    },

    async submitData() {
      let response = null
      if (this.formData.id) {
        response = await api.updateById(this.formData.id, this.formData)
      } else {
        response = await api.add(this.formData)
      }

      if (response.data === 'ok') {
        this.$message({
          message: '保存成功',
          type: 'success'
        })
        this.handleClose()
      } else {
        this.$message({
          message: '保存失败',
          type: 'error'
        })
      }
    }
  }
}
</script>
