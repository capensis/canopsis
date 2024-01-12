<template>
  <v-tooltip
    class="view-editing-btn"
    left
  >
    <template #activator="{ on }">
      <v-btn
        v-on="on"
        :input-value="editing"
        :loading="editingProcess"
        class="view-editing-btn__button"
        fab
        dark
        small
        @click.stop="toggleEditing"
      >
        <v-icon small>
          edit
        </v-icon>
        <v-icon small>
          done
        </v-icon>
      </v-btn>
    </template>
    <div>
      <div>{{ $t('common.toggleEditView') }}  (ctrl + e / command + e)</div>
      <div class="font-italic">
        {{ $t('common.toggleEditViewSubtitle') }}
      </div>
    </div>
  </v-tooltip>
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
