// components/DetailModal.js
import React from 'react';
import Modal from 'react-modal';

const DetailModal = ({ isOpen, onRequestClose, detailData }) => {
  if (!detailData) return null;

  const renderTableRows = (prefix, point, randZero, randZeroOne) => {
    if (point === 0 && randZero === 0 && randZeroOne === 0) {
      return null;
    }

    return (
      <>
        <tr>
          <th>{`${prefix} Point`}</th>
          <td>{point}</td>
        </tr>
        <tr>
          <th>{`${prefix} RandZero`}</th>
          <td>{randZero}</td>
        </tr>
        <tr>
          <th>{`${prefix} RandZeroOne`}</th>
          <td>{randZeroOne}</td>
        </tr>
      </>
    );
  };

  return (
    <Modal isOpen={isOpen} onRequestClose={onRequestClose} contentLabel="Detail Data">
      <h2>Detail Data</h2>
      <table>
        <tbody>
          {renderTableRows('Local', detailData.LocalPoint, detailData.LocalRandZero, detailData.LocalRandZeroOne)}
          {renderTableRows('EtcGG', detailData.EtcGGPoint, detailData.EtcGGRandZero, detailData.EtcGGRandZeroOne)}
          {renderTableRows('Etc', detailData.EtcPoint, detailData.EtcRandZero, detailData.EtcRandZeroOne)}
        </tbody>
      </table>
      <button onClick={onRequestClose}>Close</button>
    </Modal>
  );
};

export default DetailModal;
