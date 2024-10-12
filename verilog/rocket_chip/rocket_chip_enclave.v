module rocket_chip_enclave (
    input wire clk,
    input wire reset,
    input wire [255:0] aes_key,          // AES encryption key for decrypting instructions
    input wire [127:0] encrypted_instr,  // Encrypted instruction memory
    input wire [127:0] iv,               // AES IV for CTR mode
    input wire instruction_valid,        // Valid signal for the instruction
    input wire tamper_detected,          // Tamper detection signal
    output reg [63:0] result,            // Result of instruction execution
    output reg done                      // Instruction execution complete signal
);

    wire [127:0] decrypted_instr;
    reg start_decryption;
    wire decryption_done;

    // AES Decryption Block
    aes256_ctr aes_decrypt (
        .clk(clk),
        .reset(reset),
        .start(start_decryption),
        .key(aes_key),
        .data_in(encrypted_instr),
        .iv(iv),
        .data_out(decrypted_instr),
        .done(decryption_done)
    );

    // Rocket Chip Core for Instruction Execution
    wire [63:0] exec_result;
    wire exec_done;
    rocket_chip_core rocket_core (
        .clk(clk),
        .resetn(~reset),
        .instruction_address(64'h0000_0000),   // Address of the instruction
        .instruction_data(decrypted_instr[63:0]),  // Decrypted instruction to be executed
        .instruction_valid(instruction_valid),
        .result(exec_result),
        .done(exec_done)
    );

    always @(posedge clk or posedge reset) begin
        if (reset || tamper_detected) begin
            result <= 64'b0;  // Clear result if tamper is detected
            done <= 1'b0;
            start_decryption <= 1'b0;
        end else if (instruction_valid) begin
            start_decryption <= 1'b1;  // Start decryption process
            if (decryption_done) begin
                result <= exec_result;  // Capture execution result
                done <= exec_done;      // Mark completion
                start_decryption <= 1'b0;
            end
        end
    end
endmodule