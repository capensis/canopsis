import { DEFAULT_USER_INTERFACE_DATE_FORMAT } from '@/constants';

import { userInterfaceToForm } from '@/helpers/entities/user-interface/form';

describe('userInterfaceToForm', () => {
  it('uses the default supported date format when app-info omits it', () => {
    expect(userInterfaceToForm({}).date_format).toBe(DEFAULT_USER_INTERFACE_DATE_FORMAT);
  });

  it('keeps the API date format unchanged', () => {
    const dateFormat = '02/01/2006 (Mon) 15:04:05';

    expect(userInterfaceToForm({ date_format: dateFormat }).date_format).toBe(dateFormat);
  });
});
