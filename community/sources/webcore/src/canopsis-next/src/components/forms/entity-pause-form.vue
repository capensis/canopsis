<template>
  <v-row>
    <v-col cols="12">
      <c-description-field
        v-model="commentField"
        :label="t('common.comment', 1)"
        :max-length="FORM_MAX_STRING_LENGTH"
        name="comment"
        autofocus
      />
    </v-col>
    <v-col v-if="withInherited" cols="12">
      <v-switch
        v-model="inheritedField"
        :label="t('settings.pause.propagate')"
        color="primary"
        hide-details
      />
    </v-col>
    <v-col cols="12">
      <c-text-field
        v-model="typeField"
        :label="t('common.type')"
        :max-length="FORM_MAX_STRING_LENGTH"
        name="type"
      />
    </v-col>
  </v-row>
</template>

<script>
import { FORM_MAX_STRING_LENGTH } from '@/constants';

import { useFormModels } from '@/hooks/form';
import { useI18n } from '@/hooks/i18n';

export default {
  props: {
    withInherited: {
      type: Boolean,
      default: true,
    },
  },

  setup() {
    const { t } = useI18n();
    const { commentField, inheritedField, typeField } = useFormModels({
      commentField: 'comment',
      inheritedField: 'inherited',
      typeField: 'type',
    });

    return {
      t,
      FORM_MAX_STRING_LENGTH,
      commentField,
      inheritedField,
      typeField,
    };
  },
};
</script>
