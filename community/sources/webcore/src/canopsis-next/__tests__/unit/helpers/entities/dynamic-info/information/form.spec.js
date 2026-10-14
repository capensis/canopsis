import { describe, expect, it } from 'vitest';

import {
  dynamicInfoInformationToForm,
  formToDynamicInfoInformation,
} from '@/helpers/entities/dynamic-info/information/form';

describe('dynamic info information form helpers', () => {
  it('converts an array value to keyed form items and back', () => {
    const information = {
      name: 'tags',
      type: 0,
      value: ['first', 'second'],
    };

    const form = dynamicInfoInformationToForm(information);

    expect(form).toMatchObject({
      name: information.name,
      type: information.type,
      value: [
        { value: 'first', key: expect.any(String) },
        { value: 'second', key: expect.any(String) },
      ],
    });
    expect(formToDynamicInfoInformation(form)).toEqual(information);
  });
});
