/** @jsxImportSource @compiled/react */
import { useState, useEffect } from 'react';
import { css } from '@compiled/react';
import { sessionStorage } from '../../sessionStorage';

export const StatusPopup = () => {
  const [currentBookingRef, setCurrentBookingRef] = useState<string | null>(null);

  useEffect(() => {
    // Load initial booking reference
    const loadBookingRef = async () => {
      const bookingRef = await sessionStorage.getCurrentBookingRef();
      setCurrentBookingRef(bookingRef);
    };

    loadBookingRef();

    // Listen for changes to the booking reference
    sessionStorage.onBookingRefChange((newBookingRef) => {
      setCurrentBookingRef(newBookingRef);
    });
  }, []);

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
          <tr>
            <td css={labelCellStyle}>
              Booking Reference
            </td>
            <td css={css({
              ...valueCellStyle,
              backgroundColor: currentBookingRef ? '#e8f5e8' : '#fff3cd'
            })}>
              {currentBookingRef || 'No booking reference found'}
            </td>
          </tr>
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

const valueCellStyle = {
  padding: '8px',
  border: '1px solid #ddd',
  fontFamily: 'monospace'
};