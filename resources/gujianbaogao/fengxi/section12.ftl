<#escape x as x?xml>
<w:document xmlns:mc="http://schemas.openxmlformats.org/markup-compatibility/2006" xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main" xmlns:w14="http://schemas.microsoft.com/office/word/2010/wordml" xmlns:w15="http://schemas.microsoft.com/office/word/2012/wordml" xmlns:w16="http://schemas.microsoft.com/office/word/2018/wordml" xmlns:w16cex="http://schemas.microsoft.com/office/word/2018/wordml/cex" xmlns:w16cid="http://schemas.microsoft.com/office/word/2016/wordml/cid" xmlns:w16se="http://schemas.microsoft.com/office/word/2015/wordml/symex" xmlns:wp14="http://schemas.microsoft.com/office/word/2010/wordprocessingDrawing" mc:Ignorable="w14 w15 w16se w16cid w16 w16cex wp14">
    <w:body>
        <w:tbl>
            <w:tblPr>
                <w:tblStyle w:val="ab"/>
                <w:tblW w:type="dxa" w:w="9072"/>
                <w:jc w:val="center"/>
                <w:tblLayout w:type="fixed"/>
                <w:tblLook w:firstColumn="1" w:firstRow="1" w:lastColumn="0" w:lastRow="0" w:noHBand="0" w:noVBand="1" w:val="04A0"/>
            </w:tblPr>
            <w:tblGrid>
                <w:gridCol w:w="2198"/>
                <w:gridCol w:w="6874"/>
            </w:tblGrid>
            <#list table as row>
                <w:tr w:rsidR="00521E35" w:rsidTr="00207185">
                    <w:trPr>
                        <w:trHeight w:val="454"/>
                        <w:jc w:val="center"/>
                    </w:trPr>
                    <w:tc>
                        <w:tcPr>
                            <w:tcW w:type="dxa" w:w="2198"/>
                            <w:shd w:color="auto" w:fill="auto"  w:val="clear"/>
                            <w:vAlign w:val="center"/>
                        </w:tcPr>
                        <w:p w:rsidP="00207185" w:rsidR="00521E35" w:rsidRDefault="00521E35">
                            <w:pPr>
                                <w:widowControl/>
                                <w:spacing w:after="0" w:line="344" w:lineRule="exact"/>
                                <w:jc w:val="center"/>
                                <w:rPr>
                                    <w:rFonts w:ascii="微软雅黑" w:cs="Microsoft JhengHei" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                                    <w:b/>
                                    <w:color w:themeColor="text1" w:themeTint="D9" w:val="262626"/>
                                    <w:sz w:val="21"/>
                                    <w:szCs w:val="21"/>
                                    <w:lang w:eastAsia="zh-CN"/>
                                </w:rPr>
                            </w:pPr>
                            <w:r>
                                <w:rPr>
                                    <w:rFonts w:ascii="微软雅黑" w:cs="Microsoft JhengHei" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                                    <w:b/>
                                    <w:color w:themeColor="text1" w:themeTint="D9" w:val="262626"/>
                                    <w:sz w:val="21"/>
                                    <w:szCs w:val="21"/>
                                    <w:lang w:eastAsia="zh-CN"/>
                                </w:rPr>
                                <w:t>${row[0]}</w:t>
                            </w:r>
                        </w:p>
                    </w:tc>
                    <w:tc>
                        <w:tcPr>
                            <w:tcW w:type="dxa" w:w="6874"/>
                            <w:vAlign w:val="center"/>
                        </w:tcPr>
                        <w:p w:rsidP="00207185" w:rsidR="00521E35" w:rsidRDefault="00521E35">
                            <w:pPr>
                                <w:widowControl/>
                                <w:tabs>
                                    <w:tab w:pos="4840" w:val="left"/>
                                </w:tabs>
                                <w:spacing w:after="0" w:line="240" w:lineRule="auto"/>
                                <w:jc w:val="center"/>
                                <w:rPr>
                                    <w:rFonts w:ascii="微软雅黑" w:cs="Microsoft JhengHei" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                                    <w:sz w:val="21"/>
                                    <w:szCs w:val="21"/>
                                    <w:lang w:eastAsia="zh-CN"/>
                                </w:rPr>
                            </w:pPr>
                            <w:r w:rsidRPr="004B1D2C">
                                <w:rPr>
                                    <w:rFonts w:ascii="微软雅黑" w:cs="Microsoft JhengHei" w:eastAsia="微软雅黑" w:hAnsi="微软雅黑"/>
                                    <w:sz w:val="21"/>
                                    <w:szCs w:val="21"/>
                                    <w:lang w:eastAsia="zh-CN"/>
                                </w:rPr>
                                <w:t>${row[1]}</w:t>
                            </w:r>
                        </w:p>
                    </w:tc>
                </w:tr>
            </#list>
        </w:tbl>
        <w:p w:rsidR="00000000" w:rsidRDefault="00521E35"/>
    </w:body>
</w:document>
</#escape>
