<template lang="pug">
  file-selector(
    ref="fileSelector",
    multiple,
    hide-details,
    @change="changeFiles"
  )
    template(#activator="{ on }")
      v-tooltip(top)
        template(#activator="{ on: tooltipOn }")
          v-btn(
            v-on="{ ...on, ...tooltipOn }",
            color="indigo",
            small,
            dark,
            fab
          )
            v-icon cloud_upload
        span {{ $t('snmpRule.uploadMib') }}
</template>

<script>
import { createNamespacedHelpers } from 'vuex';

import { getFileTextContent } from '@/helpers/file/file-select';

import FileSelector from '@/components/forms/fields/file-selector.vue';

const { mapActions } = createNamespacedHelpers('snmpMib');

export default {
  components: { FileSelector },
  methods: {
    ...mapActions({
      uploadSnmpMib: 'upload',
    }),

    async changeFiles(files = []) {
      try {
        if (files.length) {
          const promises = files.sort((a, b) => {
            const matchesA = a.name.match(/^(\d+)/) || ['0'];
            const matchesB = b.name.match(/^(\d+)/) || ['0'];

            return parseInt(matchesA[0], 10) - parseInt(matchesB[0], 10);
          }).map(file => getFileTextContent(file));

          const results = await Promise.all(promises);

          const { data: [{ msg, data }] } = await this.uploadSnmpMib({
            data: results.join(''),
          });

          const popup = {
            text: msg,
            autoClose: 10000,
          };

          if (data && data.length) {
            this.$popups.error(popup);
          } else {
            this.$popups.success(popup);
          }
        }
      } catch (err) {
        this.$popups.error({
          text: this.$t('errors.default'),
        });
      }

      this.$refs.fileSelector.clear();
    },
  },
};
</script>
