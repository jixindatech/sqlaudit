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
      <el-form-item label="规则名称：" prop="name">
        <el-input v-model="formData.name" />
      </el-form-item>
      <el-form-item label="规则类型：" prop="type">
        <el-radio-group v-model="formData.type">
          <el-radio :label="1">允许类型</el-radio>
          <el-radio :label="2">拒绝类型</el-radio>
        </el-radio-group>
      </el-form-item>
      <el-form-item label="数据库用户" prop="user">
        <el-input v-model="formData.user" />
      </el-form-item>
      <el-form-item label="用户IP" prop="ip">
        <el-input v-model="formData.ip" />
      </el-form-item>
      <el-form-item label="数据库" prop="db">
        <el-input v-model="formData.db" />
      </el-form-item>
      <el-form-item label="操作类型" prop="op">
        <el-select v-model.number="formData.op">
          <el-option
            v-for="item in sqlOptions"
            :key="item.key"
            :label="item.value"
            :value="item.key"
            size="mini"
          />
        </el-select>
      </el-form-item>
      <el-form-item prop="match">
        <template slot="label">
          <span style="position:relative">
            <span>匹配方式</span>
            <el-tooltip style="position:absolute;right:-8px;" class="item" effect="dark" placement="top">
              <div slot="content">
                <p>如果指定了匹配方式：需要在输入框内输入匹配的字符串或正则</p>
              </div>
              <i class="el-icon-question table-msg" />
            </el-tooltip>
          </span>
        </template>
        <el-select v-model="formData.match" clearable placeholder="请选择">
          <el-option
            v-for="item in sqlMatchOptions"
            :key="item.key"
            :label="item.value"
            :value="item.key"
          />
        </el-select>
        <el-input v-model="formData.sql" clearable style="margin-top: 10px;" />
      </el-form-item>
      <el-form-item label="优先级：" prop="priority">
        <el-input-number v-model.number="formData.priority" :min="1" />
      </el-form-item>
      <el-form-item label="是否告警：" prop="alert">
        <el-radio-group v-model="formData.alert">
          <el-radio :label="0">否</el-radio>
          <el-radio :label="1">是</el-radio>
        </el-radio-group>
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
import api from '@/api/rule'
import { sqlOptions, sqlMatchOptions } from '@/utils/const'
import { validIP } from '@/utils/validate'

export default {
  props: {
    title: { // 弹窗的标题
      type: String,
      default: ''
    },
    visible: { // 弹出窗口，true弹出
      type: Boolean,
      default: false
    },
    formData: { // 提交表单数据
      type: Object,
      default: () => ({})
    },
    remoteClose: {
      type: Function,
      default: () => () => {}
    }
  },

  data() {
    const validateIP = (rule, value, callback) => {
      if (value !== undefined && value !== '' && !validIP(value)) {
        callback(new Error('请输入正确的IP地址'))
      } else {
        callback()
      }
    }
    const validateMatch = (rule, value, callback) => {
      if (this.formData.match === undefined || this.formData.match === '') {
        if (this.formData.sql !== undefined && this.formData.sql !== '') {
          return callback(new Error('请选择匹配方式'))
        }
      } else {
        if (this.formData.sql === undefined || this.formData.sql === '') {
          return callback(new Error('请输入匹配内容'))
        }
      }

      if (this.formData.match === '') {
        this.formData.match = undefined
        this.formData.sql = undefined
      }

      return callback()
    }
    return {
      sqlMatchOptions,
      sqlOptions,
      rules: {
        name: [ // prop值
          { required: true, message: '请输入规则名称', trigger: 'blur' }
        ],
        type: [
          { required: true, message: '请选择规则类型', trigger: 'change' }
        ],
        ip: [
          { trigger: 'change', validator: validateIP }
        ],
        op: [
          { required: true, message: '请选择操作指令', trigger: 'change' }
        ],
        match: [
          { trigger: 'change', validator: validateMatch }
        ],
        priority: [
          { required: true, message: '请输入优先级', type: 'number', trigger: 'change' }
        ],
        alert: [
          { required: true, message: '请选择是否告警', trigger: 'change' }
        ]
      }
    }
  },

  methods: {
    // 关闭窗口
    handleClose() {
      // 将表单清空
      this.$refs['formData'].resetFields()
      // 注意不可以通过  this.visible = false来关闭，因为它是父组件的属性
      this.remoteClose()
    },

    // 提交表单数据
    submitForm(formName) {
      this.$refs[formName].validate((valid) => {
        if (valid) {
          // 校验通过，提交表单数据
          this.submitData()
        } else {
          // console.log('error submit!!');
          return false
        }
      })
    },

    async submitData() {
      let response = null
      if (this.formData.id) {
        // 编辑
        response = await api.updateById(this.formData.id, this.formData)
      } else {
        // 新增
        response = await api.add(this.formData)
      }

      if (response.data === 'ok') {
        this.$message({
          message: '保存成功',
          type: 'success'
        })
        // 关闭窗口
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
