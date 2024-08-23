import React from 'react';
import './table.style.css';

const Table = ({ aptData, handleButtonClick }) => {
  return (
    <table className="tbl_st">
      <thead>
        <tr>
          <th>순서</th>
          <th>지역</th>
          <th>구분</th>            
          <th>주택명</th> 
          <th>모집공고일</th>  
          <th>특별공급 물량 및 경쟁률</th>
          <th>일반공급 물량 및 경쟁률</th>
        </tr>
      </thead>
      <tbody>
        {aptData.length > 0 ? (
          aptData.map((e, i) => (
            <tr key={e.HOUSE_MANAGE_NO}>
              <td>{i + 1}</td>
              <td>{e.SUBSCRPT_AREA_CODE_NM}</td>              
              <td>{e.HOUSE_DTL_SECD_NM}</td>
              <td>{e.HOUSE_NM}</td>
              <td>{e.RCRIT_PBLANC_DE}</td>              
              <td>
                <button type="button" onClick={() => handleButtonClick(e.DETAIL_DATA, e.HOUSE_DTL_SECD_NM, e.HOUSE_NM)}>특별 물량</button>
              </td>
              <td>
                <button type="button" onClick={() => handleButtonClick(e.DETAIL_DATA, e.HOUSE_DTL_SECD_NM, e.HOUSE_NM)}>일반 물량</button>
              </td>
            </tr>
          ))
        ) : (
          <tr><td colSpan="6">데이터가 없습니다</td></tr>
        )}
      </tbody>
    </table>
  );
};

export default Table;
