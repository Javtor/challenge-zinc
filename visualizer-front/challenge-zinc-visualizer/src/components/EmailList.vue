<template>

  <div class="grid grid-cols-2 gap-4 my-4">
    <div class="relative overflow-x-auto sm:rounded-lg">
      <table class="w-full text-sm text-left shadow-md text-gray-500 dark:text-gray-400">
        <thead class="text-xs text-gray-700 uppercase bg-gray-50 dark:bg-gray-700 dark:text-gray-400">
          <tr>
            <th scope="col" class="px-6 py-3">
              Subject
            </th>
            <th scope="col" class="px-6 py-3">
              From
            </th>
            <th scope="col" class="px-6 py-3">
              To
            </th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(email, i) in emails" :key="i" @click="selectedEmail = email" :class="{
            'bg-blue-600 rounded dark:bg-blue-500 text-white': selectedEmail === email,
            'hover:bg-gray-50 dark:hover:bg-gray-600': selectedEmail !== email
          }" class="bg-white border-b dark:bg-gray-800 dark:border-gray-700 cursor-pointer">
            <th scope="row" :class="{
              'text-white': selectedEmail === email,
              'text-gray-900 dark:text-white': selectedEmail !== email
            }" class="px-6 py-4 font-medium whitespace-nowrap">
              {{ email.subject }}
            </th>
            <td class="px-6 py-4">
              {{ email.from }}
            </td>
            <td class="px-6 py-4">
              {{ email.to.join(', ') }}
            </td>
          </tr>
        </tbody>
      </table>


      <div class="flex flex-col items-center">
        <!-- Help text -->
        <span class="text-sm text-gray-700 dark:text-gray-400 mt-2">
          Showing <span class="font-semibold text-gray-900 dark:text-white">{{ (page - 1) * perPage + 1 }}</span> to
          <span class="font-semibold text-gray-900 dark:text-white">{{
            Math.min((page *
              perPage), totalEmails)
          }}</span> of <span class="font-semibold text-gray-900 dark:text-white">{{ totalEmails }}</span> Emails
        </span>
        <span class="text-sm text-gray-700 dark:text-gray-400 mt-2">
          Page <span class="font-semibold text-gray-900 dark:text-white">{{ page }}</span> of
          <span class="font-semibold text-gray-900 dark:text-white">{{ pages }}</span>
        </span>
        <div class="inline-flex mt-1 xs:mt-0">
          <!-- Buttons -->
          <button @click="previousPage" :disabled="page === 1"
          :style="{visibility: page !== 1 ? 'visible' : 'hidden'}"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-800 rounded-l hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
            <svg aria-hidden="true" class="w-5 h-5 mr-2" fill="currentColor" viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd"
                d="M7.707 14.707a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 1.414L5.414 9H17a1 1 0 110 2H5.414l2.293 2.293a1 1 0 010 1.414z"
                clip-rule="evenodd"></path>
            </svg>
            Prev
          </button>
          <button @click="nextPage" :disabled="page === pages"
          :style="{visibility: page !== pages ? 'visible' : 'hidden'}"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-gray-800 border-0 border-l border-gray-700 rounded-r hover:bg-gray-900 dark:bg-gray-800 dark:border-gray-700 dark:text-gray-400 dark:hover:bg-gray-700 dark:hover:text-white">
            Next
            <svg aria-hidden="true" class="w-5 h-5 ml-2" fill="currentColor" viewBox="0 0 20 20"
              xmlns="http://www.w3.org/2000/svg">
              <path fill-rule="evenodd"
                d="M12.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L14.586 11H3a1 1 0 110-2h11.586l-2.293-2.293a1 1 0 010-1.414z"
                clip-rule="evenodd"></path>
            </svg>
          </button>
        </div>
      </div>
    </div>

    <div v-if="selectedEmail" class="bg-gray-900 h-screen p-4 sm:rounded-lg" >
      <h2 class="text-3xl text-gray-900 dark:text-white font-medium mb-3" style="white-space: pre-wrap;">{{ selectedEmail.subject }}</h2>
      <pre class="text-lg text-gray-500 dark:text-gray-400 overflow-auto h-full" style="white-space: pre-wrap;">{{ selectedEmail.body }}</pre>
    </div>

  </div>

</template>

<script>
export default {
  props: {
    emails: {
      type: Array,
      default: () => []
    },
    totalEmails: {
      type: Number,
      default: 0
    },
    page: {
      type: Number,
      default: 1
    },
    perPage: {
      type: Number,
      default: 10
    }
  },
  computed: {
    pages() {
      return Math.ceil(this.totalEmails / this.perPage)
    }
  },
  data() {
    return {
      selectedEmail: null
    }
  },
  methods: {
    previousPage() {
      if (this.page > 1) {
        this.$emit('update:page', this.page - 1)
      }
    },
    nextPage() {
      if (this.page < this.pages) {

        this.$emit('update:page', this.page + 1)
      }
    }
  }
}
</script>
