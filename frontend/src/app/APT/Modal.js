import React, { useState } from 'react';
import './modal.style.css';

const Modal = ({ modalData, houseDtlSecdNm, houseName, closeModal }) => {  
  return (   
        <div className="modal">
          <div className="modal-content">
            <span className="close" onClick={closeModal}>&times;</span>
            <h2 className="modal-title">{houseName}</h2>
            {modalData.map((item, index) => {
              const localData = item.LOCAL_POINT || item.LOCAL_RAND_ZERO || item.LOCAL_RAND_ZERO_ONE;
              const etcKYGData = item.ETC_KYG_POINT || item.ETC_KYG_RAND_ZERO || item.ETC_KYG_RAND_ZERO_ONE;
              const etcData = item.ETC_POINT || item.ETC_RAND_ZERO || item.ETC_RAND_ZERO_ONE;
              const kukminData =item.LOCAL_KUKMIN;
      
              return (
                <table key={index}>
                  <tbody>
                    <tr>
                      <th>주택형</th>
                      <th>가격</th>
                      <th>일반물량</th>
                      <th>지역</th>
                      <th>구분</th>
                      <th>물량</th>
                    </tr>
                    <tr>
                      <td rowSpan="0">{item.HOUSE_TY}</td>
                      <td rowSpan="0">{item.LTTOT_TOP_AMOUNT}억</td>
                      <td rowSpan="0">{item.SUPLY_HSHLDCO}</td>
                      <td rowSpan="3">해당</td>
                      {
                        kukminData > 0 ?
                        <>
                        <td>{Number(item.HOUSE_TY.replace(/[^0-9.]/g, '')) > 40 ? "高저축자" : "多납입횟수자"}</td>
                        <td>{item.LOCAL_KUKMIN}</td>
                        </>
                       : kukminData == 0  && localData> 0 ? 
                        <>
                        <td>高가점자</td>
                        <td>{item.LOCAL_POINT}</td>
                        </>
                        :
                        <>                      
                        <td></td>
                        <td>물량이 없습니다</td>                          
                        </>
                      }                 
                    </tr>
                    {localData > 0 && (
                      <>
                        <tr>
                          <td>랜덤 - 무주택</td>
                          <td>{item.LOCAL_RAND_ZERO}</td>
                        </tr>
                        <tr>
                          <td>랜덤 - 무주택,일주택</td>
                          <td>{item.LOCAL_RAND_ZERO_ONE}</td>
                        </tr>
                      </>
                    )}
                    {etcKYGData > 0 && (
                      <>
                        <tr>
                          <td rowSpan="3">해당, 기타경기</td>
                          <td>高가점자</td>
                          <td>{item.ETC_KYG_POINT}</td>
                        </tr>
                        <tr>
                          <td>랜덤 - 무주택</td>
                          <td>{item.ETC_KYG_RAND_ZERO}</td>
                        </tr>
                        <tr>
                          <td>랜덤 - 무주택,일주택</td>
                          <td>{item.ETC_KYG_RAND_ZERO_ONE}</td>
                        </tr>
                      </>
                    )}
                    {etcData > 0 && (
                      <>
                        <tr>
                          {etcKYGData > 0 && <td rowSpan="3">해당, 기타경기, 기타</td>}
                          {etcKYGData <= 0 && <td rowSpan="3">해당, 기타</td>}
                          <td>高가점자</td>
                          <td>{item.ETC_POINT}</td>
                        </tr>
                        <tr>
                          <td>랜덤 - 무주택</td>
                          <td>{item.ETC_RAND_ZERO}</td>
                        </tr>
                        <tr>
                          <td>랜덤 - 무주택,일주택</td>
                          <td>{item.ETC_RAND_ZERO_ONE}</td>
                        </tr>
                      </>
                    )}
                  </tbody>
                </table>
              );
            })}
        <div>
              <p className="highlight">
                ※ 특별공급 미달 잔여물량은 포함안함
                <br />
                ※ 할당된 물량을 모두 가져갔다고 가정했을때의 결과                               
              </p>
              {/*  ※ 투기 과열지역에 대한 계산은 미구현 */}
               {/* 만약 경쟁률 계산시 이를 고려해야함 */}
            {/* 데이터 안에 RceptEndde는 이 경쟁률인지 아닌지를 파악하는 용도로 사용됨 */}
            {/* 물량확인 버튼, 날짜 넘어가면 경쟁률 확인 버튼 */}
            </div>
      </div>
    </div>
  );
};

export default Modal;
