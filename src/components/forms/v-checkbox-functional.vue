<template lang="pug" functional>
  .v-input.v-input--selection-controls.v-input--checkbox
    .v-input__control
      .v-input__slot
        .v-input--selection-controls__input(@click="listeners.change(!props.inputValue)")
          .v-input--selection-controls__ripple.accent--text(v-ripple="{ center: true }")
          i.v-icon.material-icons(
          :class="{ 'accent--text': props.inputValue }"
          ) {{ props.inputValue ? 'check_box' : 'check_box_outline_blank' }}
        label(v-show="props.label !== ''") {{ props.label }}
</template>
