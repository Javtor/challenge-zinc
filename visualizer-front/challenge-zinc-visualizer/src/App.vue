
<template>


  <nav class="p-3 border-gray-200 rounded bg-gray-50 dark:bg-gray-800 dark:border-gray-700">
    <div class="container flex flex-wrap items-center justify-between mx-auto">
      <span class="self-center text-4xl font-semibold whitespace-nowrap dark:text-white">SWE Challenge - Zinc</span>
    </div>
  </nav>

  <div class="mx-8 h-screen">
    <form class="my-8" @submit.prevent="newSearch">
      <label for="default-search" class="mb-2 text-sm font-medium text-gray-900 sr-only dark:text-white">Search</label>
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          <svg aria-hidden="true" class="w-5 h-5 text-gray-500 dark:text-gray-400" fill="none" stroke="currentColor"
            viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </div>
        <input type="search" id="default-search" v-model="query"
          class="block w-full p-4 pl-10 text-sm text-gray-900 border border-gray-300 rounded-lg bg-gray-50 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500"
          placeholder="Search Emails..." required>
        <button type="submit"
          class="text-white absolute right-2.5 bottom-2.5 bg-blue-700 hover:bg-blue-800 focus:ring-4 focus:outline-none focus:ring-blue-300 font-medium rounded-lg text-sm px-4 py-2 dark:bg-blue-600 dark:hover:bg-blue-700 dark:focus:ring-blue-800">Search</button>
      </div>
    </form>

    <div role="status" v-if="loading" class="mx-8 flex items-center justify-center">
      <svg aria-hidden="true"
        class="object-center w-8 h-8 mr-2 text-gray-200 animate-spin dark:text-gray-600 fill-blue-600"
        viewBox="0 0 100 101" fill="none" xmlns="http://www.w3.org/2000/svg">
        <path
          d="M100 50.5908C100 78.2051 77.6142 100.591 50 100.591C22.3858 100.591 0 78.2051 0 50.5908C0 22.9766 22.3858 0.59082 50 0.59082C77.6142 0.59082 100 22.9766 100 50.5908ZM9.08144 50.5908C9.08144 73.1895 27.4013 91.5094 50 91.5094C72.5987 91.5094 90.9186 73.1895 90.9186 50.5908C90.9186 27.9921 72.5987 9.67226 50 9.67226C27.4013 9.67226 9.08144 27.9921 9.08144 50.5908Z"
          fill="currentColor" />
        <path
          d="M93.9676 39.0409C96.393 38.4038 97.8624 35.9116 97.0079 33.5539C95.2932 28.8227 92.871 24.3692 89.8167 20.348C85.8452 15.1192 80.8826 10.7238 75.2124 7.41289C69.5422 4.10194 63.2754 1.94025 56.7698 1.05124C51.7666 0.367541 46.6976 0.446843 41.7345 1.27873C39.2613 1.69328 37.813 4.19778 38.4501 6.62326C39.0873 9.04874 41.5694 10.4717 44.0505 10.1071C47.8511 9.54855 51.7191 9.52689 55.5402 10.0491C60.8642 10.7766 65.9928 12.5457 70.6331 15.2552C75.2735 17.9648 79.3347 21.5619 82.5849 25.841C84.9175 28.9121 86.7997 32.2913 88.1811 35.8758C89.083 38.2158 91.5421 39.6781 93.9676 39.0409Z"
          fill="currentFill" />
      </svg>
      <span class="sr-only">Loading...</span>
    </div>

    <div v-if="!emails.length" class="flex items-center justify-center py-12">
      <h1 class="text-4xl dark:text-white">No Results Found</h1>
    </div>
    <div v-else>
      <label for="per-page" class="block mb-2 text-base font-medium text-gray-900 dark:text-white">Items per
        page:</label>
      <select id="per-page" v-model="perPage"
        class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-blue-500 focus:border-blue-500 block p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-blue-500 dark:focus:border-blue-500">
        <option selected value="10">10</option>
        <option value="25">25</option>
        <option value="50">50</option>
      </select>
      <EmailList :emails="emails" :totalEmails="totalEmails" :page="page" :perPage="perPage" @update:page="updatePage" />
    </div>
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
      const { data } = await axios.get(`http://challenge-zinc-visualizer-back:3000/search?q=${this.query}&page=${this.page}&per_page=${this.perPage}`)
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