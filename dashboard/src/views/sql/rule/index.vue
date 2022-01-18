<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item label="规则名称：">
        <el-input v-model.trim="query.name" />
      </el-form-item>
      <el-form-item label="排序方式">
        <el-select v-model.number="query.sort">
          <el-option
            v-for="item in sortOptions"
            :key="item.key"
            :label="item.value"
            :value="item.key"
            size="mini"
          />
        </el-select>
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
      <el-table-column align="center" prop="name" label="规则名称" />
      <el-table-column align="center" prop="type" label="规则类型">
        <template slot-scope="scope">
          <el-tag :type="scope.row.type | typeFilter">
            {{ scope.row.type===1 ? '允许类型': '拒绝类型' }}
          </el-tag>
          <el-tag :type="scope.row.type | rulTypeFilter">
            {{ scope.row.ruletype===1 ? '字符串匹配': '指纹匹配' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="user" label="数据库用户" />
      <el-table-column align="center" prop="ip" label="用户IP" />
      <el-table-column align="center" prop="db" label="数据库" />
      <el-table-column align="center" prop="op" label="命令内容">
        <template slot-scope="scope">
          {{ scope.row.match === 0 ? (scope.row.op === 0 ? scope.row.sql: sqlOpMap[scope.row.op]) : sqlOpMap[scope.row.op] + ' ' + sqlMatchMap[scope.row.match] + ' ' + scope.row.sql }}
        </template>
      </el-table-column>
      <el-table-column align="center" prop="priority" label="优先级">
        <template v-if="scope.row.ruletype === 1" slot-scope="scope">
          {{ scope.row.priority }}
        </template>
      </el-table-column>
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
import api from '@/api/rule'
import Edit from './edit'
import { sqlOpMap, sortOptions, sqlMatchMap } from '@/utils/const'

export default {
  components: { Edit },
  filters: {
    typeFilter(status) {
      const statusMap = { 1: 'success', 2: 'danger' }
      return statusMap[status]
    },
    rulTypeFilter(status) {
      const statusMap = { 1: 'primary', 2: 'info' }
      return statusMap[status]
    }
  },

  data() {
    return {
      sortOptions,
      sqlOpMap,
      sqlMatchMap,
      list: [],
      listLoading: true,
      page: { // 分页对象
        current: 1, // 当前页码
        size: 20, // 每页显示多少条
        total: 0 // 总记录数
      },
      query: { sort: 1 },
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
