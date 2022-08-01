<template lang="pug">
  v-tooltip(left)
    v-btn(
      slot="activator",
      :input-value="editing",
      :loading="editingProcess",
      fab,
      dark,
      small,
      @click.stop="toggleEditing"
    )
      v-icon edit
      v-icon done
    div
      div {{ $t('common.toggleEditView') }}  (ctrl + e / command + e)
      div.font-italic {{ $t('common.toggleEditViewSubtitle') }}
</template>

<script>
import { activeViewMixin } from '@/mixins/active-view';

export default {
  mixins: [activeViewMixin],
};
</script>
