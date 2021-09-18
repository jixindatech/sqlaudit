<template>
  <div class="app-container">
    <!-- 条件查询 -->
    <el-form :inline="true" :model="query" size="mini">
      <el-form-item label="分类名称：">
        <el-input v-model.trim="query.name" />
      </el-form-item>
      <el-form-item label="状态：">
        <el-select v-model="query.status" clearable filterable style="width: 85px">
          <el-option
            v-for="item in statusOptions"
            :key="item.code"
            :label="item.name"
            :value="item.code"
          />
        </el-select>
      </el-form-item>
      <el-form-item>
        <el-button icon="el-icon-search" type="primary" @click="queryData">查询</el-button>
        <el-button icon="el-icon-refresh" @click="reload">重置</el-button>
        <el-button icon="el-icon-circle-plus-outline" type="primary" @click="openAdd">新增</el-button>
      </el-form-item>
    </el-form>

    <el-table
      v-loading="listLoading"
      :data="list"
      element-loading-text="Loading"
      border
      fit
      highlight-current-row
    >
      <el-table-column align="center" label="ID" width="95">
        <template slot-scope="scope">
          {{ scope.$index }}
        </template>
      </el-table-column>
      <el-table-column label="Title">
        <template slot-scope="scope">
          {{ scope.row.title }}
        </template>
      </el-table-column>
      <el-table-column label="Author" width="110" align="center">
        <template slot-scope="scope">
          <span>{{ scope.row.author }}</span>
        </template>
      </el-table-column>
      <el-table-column label="Pageviews" width="110" align="center">
        <template slot-scope="scope">
          {{ scope.row.pageviews }}
        </template>
      </el-table-column>
      <el-table-column class-name="status-col" label="Status" width="110" align="center">
        <template slot-scope="scope">
          <el-tag :type="scope.row.status | statusFilter">{{ scope.row.status }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column align="center" prop="created_at" label="Display_time" width="200">
        <template slot-scope="scope">
          <i class="el-icon-time" />
          <span>{{ scope.row.display_time }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script>
import { getList } from '@/api/table'
const statusOptions = [
  { code: 0, name: '禁用' },
  { code: 1, name: '正常' }
]

export default {
  filters: {
    statusFilter(status) {
      // 0 是禁用， 1是正常
      const statusMap = { 0: 'danger', 1: 'success' }
      return statusMap[status]
    }
  },
  data() {
    return {
      list: null,
      listLoading: true,
      page: { // 分页对象
        current: 1, // 当前页码
        size: 20, // 每页显示多少条
        total: 0 // 总记录数
      },
      query: {},
      statusOptions
    }
  },

  created() {
    this.fetchData()
  },
  methods: {
    fetchData() {
      this.listLoading = true
      getList().then(response => {
        this.list = response.data.items
        this.listLoading = false
      })
    },
    queryData() {
      // 将页码变为1，第1页
      this.page.current = 1
      this.fetchData()
    },
    reload() {
      this.query = {}
      this.fetchData()
    }

  }
}
</script>
