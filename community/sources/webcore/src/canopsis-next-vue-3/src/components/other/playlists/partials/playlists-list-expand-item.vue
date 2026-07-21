<template>
  <div class="pa-4">
    <v-table density="compact">
      <thead>
        <tr>
          <th>{{ t('common.group') }}</th>
          <th>{{ t('common.view') }}</th>
          <th>{{ t('common.tab') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="tabId in playlist.tabs_list || []" :key="tabId">
          <td colspan="3">
            {{ tabId }}
          </td>
        </tr>
      </tbody>
    </v-table>
  </div>
</template>

<script>
import { useI18n } from '@/hooks/i18n'

export default {
  props: {
    playlist: {
      type: Object,
      required: true,
    },
  },

  setup() {
    const { t } = useI18n()

    return { t }
  },
}
</script>
