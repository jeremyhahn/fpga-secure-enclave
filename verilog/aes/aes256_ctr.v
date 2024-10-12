module aes256_ctr (
    input wire clk,
    input wire rst,
    input wire [127:0] plaintext,
    input wire [255:0] key,
    output wire [127:0] ciphertext,
    output wire done
);

    wire ready;
    reg init;
    reg next;

    // AES core instantiation from secworks library
    aes_core aes_inst (
        .clk(clk),
        .reset_n(~rst),
        .init(init),
        .next(next),
        .key(key),
        .keylen(1'b1),         // 256-bit AES
        .block(plaintext),
        .ready(ready),
        .result(ciphertext),
        .done(done)
    );

    // Control logic for AES operation
    always @(posedge clk or posedge rst) begin
        if (rst) begin
            init <= 1'b1;
            next <= 1'b0;
        end else if (ready && !done) begin
            next <= 1'b1;
        end else begin
            next <= 1'b0;
        end
    end

endmodule