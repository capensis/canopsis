<template>
  <c-layout column>
    <c-layout class="icon-form__title-row" align="center" gap="3">
      <div v-if="iconOriginalTitle" class="icon-form__preview">
        <custom-icon :content="content" />
      </div>

      <c-text-field
        v-model="title"
        :label="$t('common.title')"
        :regex="titleRegex"
        :regex-message="titleRegexMessage"
        name="title"
        required
        autofocus
      />
    </c-layout>

    <c-enabled-field v-model="fillBorder" :label="$t('icon.fillBorder')" />

    <icon-file-selector
      v-if="!iconOriginalTitle"
      v-model:file="file"
      :max-file-size="maxFileSize"
      class="mt-2"
    />
  </c-layout>
</template>

<script>
import { ref } from 'vue';

import { ICON_TITLE_REGEX, MAX_ICON_SIZE_IN_KB } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

import CustomIcon from '../partials/custom-icon.vue';

import IconFileSelector from './fields/icon-file-selector.vue';

export default {
  components: {
    CustomIcon,
    IconFileSelector,
  },

  setup() {
    const { t } = useI18n();
    const { title, fillBorder, content, file } = useFormModels({
      title: 'title',
      fillBorder: 'fill_border',
      content: 'content',
      file: 'file',
    });

    const iconOriginalTitle = ref(title.value);
    const titleRegex = ICON_TITLE_REGEX;
    const maxFileSize = MAX_ICON_SIZE_IN_KB;
    const titleRegexMessage = t('icon.validation.iconTitleRegex');

    return {
      title,
      fillBorder,
      content,
      file,
      iconOriginalTitle,
      titleRegex,
      titleRegexMessage,
      maxFileSize,
    };
  },
};
</script>

<style lang="scss">
.icon-form {
  &__title-row {
    min-width: 0;
  }

  &__preview {
    flex: 0 0 48px;
  }
}
</style>
