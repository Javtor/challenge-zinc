<template>
  <div class="grid grid-cols-2 gap-4">
    <div>
      <table class="table-auto w-full">
        <thead>
          <tr>
            <th class="px-4 py-2">Subject</th>
            <th class="px-4 py-2">From</th>
            <th class="px-4 py-2">To</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="(email, i) in emails" :key="i" @click="selectedEmail = email" class="hover:bg-gray-200"
            :class="{ 'bg-indigo-500 text-white': selectedEmail === email }">
            <td class="border px-4 py-2">{{ email.subject }}</td>
            <td class="border px-4 py-2">{{ email.from }}</td>
            <td class="border px-4 py-2">{{ email.to.join(', ') }}</td>
          </tr>
        </tbody>
      </table>
      <div v-if="!emails.length">No Results Found</div>
      <div class="text-center my-4">
        Page {{ page }} of {{ pages }}
      </div>
      <div class="my-4">
        <button class="button p-2" @click="previousPage" :disabled="page === 1">Previous</button>
        <button class="button p-2" @click="nextPage" :disabled="page === pages">Next</button>
      </div>
      
    </div>

    <div v-if="selectedEmail" class="bg-gray-200">
      <h2 class="p-4 text-2xl">{{ selectedEmail.subject }}</h2>
      <pre class="overflow-x-auto p-4 rounded-md" style="white-space: pre-wrap;">{{ selectedEmail.body }}</pre>
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
