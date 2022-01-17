<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item label="规则名称：">
        <el-input v-model.trim="query.name" />
      </el-form-item>
      <el-form-item>
        <el-button icon="el-icon-search" type="primary" @click="queryData">查询</el-button>
        <el-button icon="el-icon-refresh" @click="reload">重置</el-button>
        <el-button icon="el-icon-circle-plus-outline" type="primary" @click="openAdd">新增</el-button>
      </el-form-item>
    </el-form>
    <edit :title="edit.title" :visible="edit.visible" :form-data="edit.formData" :remote-close="remoteClose" />
    <el-table
      :data="list"
      stripe
      border
      style="width: 100%"
    >
      <el-table-column align="center" type="index" label="序号" width="60" />
      <el-table-column align="center" prop="name" label="名称" />
      <el-table-column align="center" prop="remark" label="备注" />
      <el-table-column align="center" label="操作" width="330">
        <template slot-scope="scope">
          <el-button
            type="success"
            size="mini"
            @click="handleEdit(scope.row.id)"
          >编辑</el-button>
          <el-button
            type="danger"
            size="mini"
            @click="handleDelete(scope.row.id)"
          >删除</el-button>
        </template>
      </el-table-column>

    </el-table>
    <el-pagination
      :current-page="page.current"
      :page-sizes="[10, 20, 50]"
      :page-size="page.size"
      layout="total, sizes, prev, pager, next, jumper"
      :total="page.total"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>

</template>

<script>
import api from '@/api/fingerprint'
import Edit from './edit'

export default {
  components: { Edit },
  data() {
    return {
      list: [],
      listLoading: true,
      page: {
        current: 1,
        size: 20,
        total: 0
      },
      query: {},
      edit: {
        title: '',
        visible: false,
        formData: {}
      }
    }
  },
  created() {
    this.fetchData()
  },
  methods: {
    async fetchData() {
      const { data } = await api.getList(this.query, this.page.current, this.page.size)
      this.list = data.data
      this.page.total = data.count
    },

    queryData() {
      // 将页码变为1，第1页
      this.page.current = 1
      this.fetchData()
    },
    // 重置
    reload() {
      this.query = {}
      this.fetchData()
    },
    remoteClose() {
      this.edit.formData = {}
      this.edit.visible = false
      this.fetchData()
    },
    // 打开新增窗口
    openAdd() {
      this.edit.visible = true
      this.edit.title = '新增'
    },
    handleSizeChange(val) {
      this.page.size = val
      this.fetchData()
    },
    handleCurrentChange(val) {
      this.page.current = val
      this.fetchData()
    },
    handleEdit(id) {
      api.getById(id).then((response) => {
        if (response.data) {
          this.edit.formData = response.data
          if (response.data.match === 0) {
            this.edit.formData.match = undefined
            this.edit.formData.sql = undefined
          }
          this.edit.title = '编辑'
          this.edit.visible = true
        }
      })
    },
    handleDelete(id) {
      this.$confirm('确认删除这条记录吗?', '提示', {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'warning'
      })
        .then(() => {
          api.deleteById(id).then((response) => {
            this.$message({
              type: response.data === 'ok' ? 'success' : 'error',
              message: '删除成功!'
            })
            this.fetchData()
          })
        })
        .catch(() => {
        })
    }
  }
}
</script>
