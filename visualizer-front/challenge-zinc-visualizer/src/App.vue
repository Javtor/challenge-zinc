
<template>
  <div class="mx-auto w-full p-6">
    <h1 class="text-center text-4xl">My Emails</h1>
    <form @submit.prevent="newSearch">
      <input type="text" v-model="query" placeholder="Search..." class="bg-white rounded-md py-2 px-4" />
      <button type="submit" class="bg-blue-500 rounded-md py-2 px-4 text-white">Search</button>
      <button type="button" class="bg-red-500 rounded-md py-2 px-4 text-white" @click="clearSearch">Clear</button>
    </form>
    <div v-if="loading" class="text-center my-4">
      <div class="spinner"></div>
    </div>
    <div v-if="emails.length" class="my-4">
      <label for="per-page" class="block">Items per page:</label>
      <select id="per-page" class="bg-white rounded-md py-2 px-4" v-model="perPage">
        <option value="10">10</option>
        <option value="25">25</option>
        <option value="50">50</option>
      </select>
    </div>
    <div v-if="emails.length" class="my-4 text-gray-700">
      Showing {{ (page - 1) * perPage + 1 }} to {{ Math.min((page * perPage), totalEmails) }} of {{ totalEmails }}
      entries.
    </div>
    <EmailList :emails="emails" :totalEmails="totalEmails" :page="page" :perPage="perPage" @update:page="updatePage" />
  </div>
</template>

<script>
import axios from 'axios'
import EmailList from './components/EmailList.vue'

export default {
  components: {
    EmailList
  },
  data() {
    return {
      query: '',
      emails: [],
      totalEmails: 0,
      page: 1,
      perPage: 10,
      loading: false
    }
  },
  methods: {
    newSearch() {
      this.page = 1
      this.searchEmails()
    },
    async searchEmails() {
      this.loading = true
      const { data } = await axios.get(`http://localhost:3000/search?q=${this.query}&page=${this.page}&per_page=${this.perPage}`)
      this.emails = data.emails
      this.totalEmails = data.total
      this.loading = false
    },
    updatePage(newPage) {
      this.page = newPage
      this.searchEmails()
      console.log(this.page)
    },
    clearSearch() {
      this.query = ''
      this.emails = []
      this.totalEmails = 0
    }
  },
  watch: {
    perPage(newValue) {
      this.page = 1
      this.searchEmails()
    }
  }
}
</script>