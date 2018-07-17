<template lang="pug">
  v-list-group
    v-list-tile(slot="activator") {{$t('settings.periodicRefresh')}}
    v-container
      v-layout
        v-flex
          v-switch(
          v-model="isEnabled",
          color="green darken-3",
          hide-details
          )
        v-flex
          v-text-field.pt-0(
          hide-details,
          type="number",
          :disabled="!isEnabled"
          )
</template>

<script>
export default {
  data() {
    return {
      isEnabled: false,
    };
  },
};
</script>
