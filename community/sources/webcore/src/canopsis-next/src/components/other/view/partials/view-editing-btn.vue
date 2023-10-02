<template lang="pug">
  v-tooltip.view-editing-btn(left)
    v-btn.view-editing-btn__button(
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

<style lang="scss">
.view-editing-btn {
  &__button  {
    border-color: #212121 !important;
    background-color: #212121 !important;

    .theme--dark & {
      border-color: #979797 !important;
      background-color: #979797 !important;
    }
  }
}
</style>
