<template>
  <file-selector
    ref="fileSelector"
    multiple
    hide-details
    @change="changeFiles"
  >
    <template #activator="{ on, attrs }">
      <v-tooltip location="top">
        <template #activator="{ props }">
          <v-btn
            v-bind="{ ...attrs, ...props }"
            :loading="pending"
            color="indigo"
            size="small"
            class="snmp-rules-list-upload-button__fab"
            icon
            @click="on.click"
          >
            <v-icon size="small">
              cloud_upload
            </v-icon>
          </v-btn>
        </template>
        <span>{{ t('snmpRule.uploadMib') }}</span>
      </v-tooltip>
    </template>
  </file-selector>
</template>

<script>
import { ref } from 'vue';

import { getFileTextContent } from '@/helpers/file/file-select';

import { useI18n } from '@/hooks/i18n';
import { usePopups } from '@/hooks/popups';
import { useSnmpMib } from '@/hooks/store/modules/snmp-mib';

import FileSelector from '@/components/forms/fields/file-selector.vue';

export default {
  components: { FileSelector },

  setup() {
    const pending = ref(false);
    const fileSelector = ref(null);
    const { t } = useI18n();
    const popups = usePopups();
    const { uploadSnmpMib } = useSnmpMib();

    const sortFiles = files => [...files].sort((a, b) => {
        const matchesA = a.name.match(/^(\d+)/) || ['0'];
        const matchesB = b.name.match(/^(\d+)/) || ['0'];

        return parseInt(matchesA[0], 10) - parseInt(matchesB[0], 10);
      });

    const changeFiles = async (files = []) => {
      if (!files.length) {
        return;
      }

      try {
        pending.value = true;

        const sortedFiles = sortFiles(files);
        const results = await Promise.all(sortedFiles.map(file => getFileTextContent(file)));

        const fileContent = results.map((content, index) => ({
          filename: sortedFiles[index].name,
          data: content,
        }));

        const { counts } = await uploadSnmpMib({
          data: { filecontent: fileContent },
        });

        popups.success({
          text: t('snmpRule.uploadedMibPopup', sortedFiles.length, { named: counts }),
          autoClose: 10000,
        });
      } catch (err) {
        console.warn(err);

        popups.error({
          text: err.error ?? t('errors.default'),
          autoClose: 10000,
        });
      } finally {
        pending.value = false;
        fileSelector.value?.clear();
      }
    };

    return {
      t,
      pending,
      fileSelector,
      changeFiles,
    };
  },
};
</script>

<style lang="scss" scoped>
.snmp-rules-list-upload-button__fab {
  --v-btn-height: 40px;

  width: 40px;
  height: 40px;
  border-radius: 50%;
}
</style>
