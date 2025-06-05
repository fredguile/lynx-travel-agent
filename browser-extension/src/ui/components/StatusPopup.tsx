/** @jsxImportSource @compiled/react */
import { useState, useEffect } from 'react';
import { css } from '@compiled/react';
import { sessionStorage } from '../../sessionStorage';

const COLORS = {
  SUCCESS_BACKGROUND: '#e8f5e8',
  WARNING_BACKGROUND: '#fff3cd'
};

export const StatusPopup = () => {
  const [currentUrl, setCurrentUrl] = useState<string | null>(null);
  const [currentBookingRef, setCurrentBookingRef] = useState<string | null>(null);

  useEffect(() => {
    // Load initial values
    const loadInitialValues = async () => {
      const url = await sessionStorage.currentUrl.get();
      const bookingRef = await sessionStorage.currentBookingRef.get();
      setCurrentUrl(url);
      setCurrentBookingRef(bookingRef);
    };

    loadInitialValues();

    // Listen for changes to the current URL
    sessionStorage.currentUrl.onPropertyChange((newUrl: string | null) => {
      setCurrentUrl(newUrl);
    });

    // Listen for changes to the booking reference
    sessionStorage.currentBookingRef.onPropertyChange((newBookingRef: string | null) => {
      setCurrentBookingRef(newBookingRef);
    });
  }, []);

  const renderTableRow = (label: string, value: string | null, defaultMessage: string) => (
    <tr>
      <td css={labelCellStyle}>
        {label}
      </td>
      <td css={valueCellStyle} style={{
        backgroundColor: value ? COLORS.SUCCESS_BACKGROUND : COLORS.WARNING_BACKGROUND
      }}>
        {value || defaultMessage}
      </td>
    </tr>
  );

  return (
    <div css={containerStyle}>
      <table css={tableStyle}>
        <tbody>
          <tr>
            <td css={headerCellStyle}>
              Property
            </td>
            <td css={headerCellLastStyle}>
              Value
            </td>
          </tr>
          {renderTableRow('Current URL', currentUrl, 'No URL found')}
          {renderTableRow('Booking Reference', currentBookingRef, 'No booking reference found')}
        </tbody>
      </table>
    </div>
  );
};

const containerStyle = css({
  padding: '16px',
  minWidth: '250px'
});

const tableStyle = css({
  width: '100%',
  borderCollapse: 'collapse',
  marginTop: '12px'
});

const headerCellStyle = css({
  padding: '8px',
  border: '1px solid #ddd',
  fontWeight: 'bold',
  backgroundColor: '#f8f9fa',
  width: '40%'
});

const headerCellLastStyle = css({
  padding: '8px',
  border: '1px solid #ddd',
  fontWeight: 'bold',
  backgroundColor: '#f8f9fa'
});

const labelCellStyle = css({
  padding: '8px',
  border: '1px solid #ddd',
  fontWeight: 'bold'
});

const valueCellStyle = css({
  padding: '8px',
  border: '1px solid #ddd',
  fontFamily: 'monospace'
});